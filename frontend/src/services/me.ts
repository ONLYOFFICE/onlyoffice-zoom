import axios from "axios";
import axiosRetry from "axios-retry";
import zoomSdk from "@zoom/appssdk";

import { MeResponse } from "src/types/user";

export const getMe = async (signal: AbortSignal | undefined = undefined) => {
  const zctx = await zoomSdk.getAppContext();
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  axiosRetry(client, {
    retries: 2,
    retryCondition: (error) => error.status !== 200,
  });
  const res = await client<MeResponse>({
    method: "GET",
    url: "/api/me",
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });

  return { response: res.data };
};
