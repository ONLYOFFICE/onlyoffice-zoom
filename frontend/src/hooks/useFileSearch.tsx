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
    staleTime: 2000,
    cacheTime: 2500,
    refetchInterval: 2000,
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
