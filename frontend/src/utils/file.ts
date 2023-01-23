import docx from "@assets/docx.svg";
import pptx from "@assets/pptx.svg";
import xls from "@assets/xls.svg";
import notsupported from "@assets/nofile.svg";

const DOCUMENT_EXTS = [
  "doc",
  "docx",
  "docm",
  "dot",
  "dotx",
  "dotm",
  "odt",
  "fodt",
  "ott",
  "rtf",
  "txt",
  "html",
  "htm",
  "mht",
  "xml",
  "pdf",
  "djvu",
  "fb2",
  "epub",
  "xps",
  "oxps",
];
const SPREADSHEET_EXTS = [
  "xls",
  "xlsx",
  "xlsm",
  "xlt",
  "xltx",
  "xltm",
  "ods",
  "fods",
  "ots",
  "csv",
];
const PRESENTATION_EXTS = [
  "pps",
  "ppsx",
  "ppsm",
  "ppt",
  "pptx",
  "pptm",
  "pot",
  "potx",
  "potm",
  "odp",
  "fodp",
  "otp",
];

const EDITABLE_EXTS = ["docx", "pptx", "xlsx"];
const OPENABLE_EXTS =
  DOCUMENT_EXTS.concat(SPREADSHEET_EXTS).concat(PRESENTATION_EXTS);

const WORD = "word";
const SLIDE = "slide";
const CELL = "cell";

const getFileExt = (filename: string): string =>
  filename.split(".").pop() || "";

export const isFileEditable = (filename: string) => {
  const ext = getFileExt(filename).toLowerCase();
  return EDITABLE_EXTS.includes(ext);
};

export const isFileSupported = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();
  return OPENABLE_EXTS.includes(e);
};

export const getFileType = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  if (DOCUMENT_EXTS.includes(e)) return WORD;
  if (SPREADSHEET_EXTS.includes(e)) return CELL;
  if (PRESENTATION_EXTS.includes(e)) return SLIDE;

  return null;
};

export const getFileIcon = (filename: string) => {
  const e = getFileExt(filename).toLowerCase();

  if (DOCUMENT_EXTS.includes(e)) return docx;
  if (SPREADSHEET_EXTS.includes(e)) return xls;
  if (PRESENTATION_EXTS.includes(e)) return pptx;

  return notsupported;
};

// TODO: Set proper defaults
export const getCreateFileUrl = (
  fileType: "docx" | "pptx" | "xlsx" | undefined
) => {
  switch (fileType) {
    case "docx":
      return encodeURIComponent(process.env.WORD_FILE || "");
    case "pptx":
      return encodeURIComponent(process.env.SLIDE_FILE || "");
    case "xlsx":
      return encodeURIComponent(process.env.SPREADSHEET_FILE || "");
    default:
      return encodeURIComponent(process.env.WORD_FILE || "");
  }
};

export const formatBytes = (bytes: number, decimals = 2) => {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
};
