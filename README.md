# AI-Driven Distributed File Orchestrator

A high-performance, multi-service system designed to categorize and organize unstructured local data using LLM-based zero-shot classification. This project leverages **Go** for concurrent I/O operations and **Python** for specialized AI inference, communicating via a structured REST IPC.

---

## 🏗 System Architecture

The system is decoupled into three functional units to optimize for both computational efficiency and I/O throughput.

### 1. Unit 1: The Dispatcher (Go)
* **Directory Traversal:** Recursively scans target paths with robust permission and existence validation.
* **MIME Detection:** Utilizes the `mimetype` library to filter and identify file types before processing.
* **Payload Preparation:** Aggregates file paths into optimized batches for transmission.

### 2. Unit 2: The Neural Inference Engine (Python/FastAPI)
* **Zero-Shot Semantic Classification:** Leverages a `BART-large-mnli` Transformer architecture utilizing **Natural Language Inference (NLI)**. By framing categorization as a premise-entailment problem, the system achieves high-accuracy classification across non-deterministic labels without requiring fine-tuned weights.
* **Context Window Management:** Implements a sliding-window text extraction strategy. To optimize latency, the system performs **Truncation & Normalization** on raw document buffers, capping input at 3,000 characters to stay within the model's effective attention span.
* **Softmax-normalized Confidence Filtering:** Computes a **Logit Distribution** across target classes. A configurable threshold ($\\\\text{score} \\\\geq 0.6$) is applied to mitigate model hallucinations and ensure deterministic file routing.
* **Asynchronous IO (asyncio):** Architected as a non-blocking service to handle concurrent requests from the Go dispatcher, maximizing throughput during large-scale ingestion.

### 3. Unit 3: The Executioner (Go)
* **Threshold Logic:** Evaluates AI confidence scores against a configurable default ($\\\\text{threshold} \\\\geq 0.6$).
* **Atomic File Operations:** Executes file movements to categorized directories or a `misc/` fallback for low-confidence results.

---

## 📡 Data Flow & IPC

The services interact through a synchronous REST-based bridge designed for resilience:

1. **Request:** Go initiates an HTTP POST to `/process-batch` with a JSON payload of discovered paths.
2. **Processing:** FastAPI handles text extraction and model inference.
3. **Response:** The results containing category labels and confidence scores are returned as an HTTP response.
4. **Action:** The Go service parses the response and triggers the `mover` package logic.

---

## 🚀 Technical Highlights

### Backend & Distributed Systems
* **Concurrency:** Go utilizes Goroutines to manage the HTTP client lifecycle and file system operations.
* **Resilience:** Implemented 30-second timeouts, custom HTTP client wrappers, and rigorous JSON unmarshaling validation.
* **Modular Design:**
    * `classifier/`: File scanning and metadata detection.
    * `network/`: HTTP/IPC resilience layer.
    * `mover/`: File system state management.

### AI Implementation
* **Model:** Zero-shot classification (BART-Large-MNLI).
* **Decisions:** Confidence-based filtering to prevent data misplacement (hallucination mitigation).
* **Parsing:** Error-tolerant UTF-8 text extraction and PDF buffer management.

---

## 🛠 Tech Stack

| Component | Technology |
| :--- | :--- |
| **Languages** | Golang, Python 3.10+ |
| **AI Framework** | BART-large-mnli, Hugging Face Transformers, PyTorch |
| **Backend** | FastAPI (Python), Standard Library (Go) |
| **Reliability** | Docker, Redis (Optional for scaling), Goroutines |

---

## ⚙️ Configuration
| Variable | Default | Description |
| :--- | :--- | :--- |
| `CONFIDENCE_THRESHOLD` | `0.6` | Minimum score required for categorization. |
| `MAX_TEXT_CHARS` | `3000` | Character limit for LLM context window. |
| `TIMEOUT_SECONDS` | `30` | Maximum wait time for IPC response. |
