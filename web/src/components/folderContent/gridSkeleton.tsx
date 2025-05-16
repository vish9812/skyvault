function GridSkeleton() {
  return (
    <div class="p-4 flex flex-wrap gap-4">
      {Array.from({ length: 10 }).map(() => (
        <div class="w-40 h-40 md:w-48 md:h-48 bg-white rounded-lg border border-border shadow-sm">
          <div class="h-28 md:h-34 rounded-t-lg border-b border-border bg-bg-subtle flex-center">
            <div class="w-14 h-14 bg-gray-200 rounded-full animate-pulse"></div>
          </div>
          <div class="p-2 flex flex-col">
            <div class="h-4 bg-gray-200 rounded w-3/4 animate-pulse"></div>
            <div class="h-3 bg-gray-200 rounded w-1/3 mt-2 animate-pulse"></div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default GridSkeleton;
