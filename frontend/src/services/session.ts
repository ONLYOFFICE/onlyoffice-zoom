import axios from "axios";
import zoomSdk from "@zoom/appssdk";

export const removeSession = async (
  signal: AbortSignal | undefined = undefined
) => {
  const zctx = await zoomSdk.getAppContext();
  const res = await axios({
    method: "DELETE",
    url: `${process.env.BACKEND_GATEWAY}/api/session`,
    headers: {
      "Content-Type": "application/json",
      "X-Zoom-App-Context": zctx.context,
    },
    signal,
  });

  return { response: res.status === 200 };
};
