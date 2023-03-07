import axios from "axios";
import zoomSdk from "@zoom/appssdk";

import { ConfigResponse } from "src/types/config";

export const fetchConfig = async (
  name: string,
  url: string,
  lang: string,
  signal?: AbortSignal
) => {
  const zctx = await zoomSdk.getAppContext();
  const res = await axios<ConfigResponse>({
    method: "GET",
    url: `${process.env.BACKEND_GATEWAY}/api/config`,
    params: {
      file_name: name,
      file_url: url,
      lang,
    },
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });
  return res.data;
};
