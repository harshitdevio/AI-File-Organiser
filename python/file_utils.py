import PyPDF2
from typing import Optional

def extract_text_from_pdf(filepath: str, max_chars: int = 3000) -> str:
    try:
        with open(filepath, 'rb') as f:
            reader = PyPDF2.PdfReader(f)
            text = ""
            for page_num in range(min(len(reader.pages), 5)):
                text += reader.pages[page_num].extract_text()
                if len(text) >= max_chars:
                    break
            return text[:max_chars]
    except Exception as e:
        raise Exception(f"PDF extraction failed: {str(e)}")

def read_file_content(filepath: str, mime_type: Optional[str] = None) -> str:
    if mime_type and 'pdf' in mime_type.lower():
        return extract_text_from_pdf(filepath)
    elif filepath.lower().endswith('.pdf'):
        return extract_text_from_pdf(filepath)
    else:
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            return f.read(3000)