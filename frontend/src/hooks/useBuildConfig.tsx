import { useQuery } from "react-query";

import { fetchConfig } from "@services/config";

export function useBuildConfig(
  id: string,
  name: string,
  url = "https://onlyoffice.com"
) {
  const { isLoading, error, data } = useQuery({
    queryKey: ["config", id, name, url],
    queryFn: ({ signal }) => fetchConfig(name, url, signal),
    staleTime: 0,
    cacheTime: 0,
    refetchOnWindowFocus: false,
  });
  return { isLoading, error, data };
}
