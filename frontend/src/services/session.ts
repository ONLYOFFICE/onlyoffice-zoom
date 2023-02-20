import axios from "axios";
import axiosRetry from "axios-retry";
import zoomSdk from "@zoom/appssdk";

export const removeSession = async (
  signal: AbortSignal | undefined = undefined
) => {
  const zctx = await zoomSdk.getAppContext();
  const client = axios.create({ baseURL: process.env.BACKEND_GATEWAY });
  axiosRetry(client, {
    retries: 2,
    retryCondition: (error) => error.status !== 200,
  });
  const res = await client({
    method: "DELETE",
    url: "/api/session",
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });

  return { response: res.status === 200 };
};
