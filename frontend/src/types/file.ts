export type File = {
  download_url: string;
  file_id: string;
  file_name: string;
  file_size: string;
  timestamp: number;
};

export type FileResponse = {
  messages: File[];
  next_page_token: string;
  page_size: number;
};
