import axios from "axios";
import zoomSdk from "@zoom/appssdk";

import { FileResponse } from "src/types/file";

export const fetchFiles = async (
  query = "",
  pageToken = "",
  signal: AbortSignal | undefined = undefined
) => {
  const zctx = await zoomSdk.getAppContext();
  const res = await axios<FileResponse>({
    method: "GET",
    url: `${process.env.BACKEND_GATEWAY}/api/files`,
    params: {
      page_size: 50,
      ...(query && { search_key: query }),
      ...(pageToken && { next_page_token: pageToken }),
    },
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });

  return res.data;
};
