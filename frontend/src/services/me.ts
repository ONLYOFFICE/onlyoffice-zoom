import axios from "axios";
import zoomSdk from "@zoom/appssdk";

export const getMe = async (signal: AbortSignal | undefined = undefined) => {
  const zctx = await zoomSdk.getAppContext();
  const res = await axios({
    method: "GET",
    url: `${process.env.BACKEND_GATEWAY}/api/me`,
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });

  return { response: res.data };
};
