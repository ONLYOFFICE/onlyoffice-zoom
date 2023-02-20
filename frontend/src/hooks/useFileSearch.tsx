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
    staleTime: 4000,
    cacheTime: 4500,
    refetchInterval: 4000,
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
