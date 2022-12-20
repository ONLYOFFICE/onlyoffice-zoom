import { useInfiniteQuery } from "react-query";

import { fetchFiles } from "@services/file";

export function useFileSearch(query = "") {
  const {
    data,
    isLoading,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: ["filesData", query],
    queryFn: ({ pageParam = "", signal }) =>
      fetchFiles(query, pageParam, signal),
    getNextPageParam: (lastPage) =>
      lastPage.nextPage ? lastPage.nextPage : undefined,
    staleTime: parseInt(process.env.FILE_STALE_TIME || "30000", 10),
    cacheTime: parseInt(process.env.FILE_CACHE_TIME || "35000", 10),
    refetchOnWindowFocus: true,
  });

  return {
    files: data?.pages
      .map((page) => page.response)
      .filter(Boolean)
      .flat(),
    isLoading,
    error,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  };
}
