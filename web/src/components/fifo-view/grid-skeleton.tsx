function GridSkeleton() {
  return (
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {Array.from({ length: 10 }).map(() => (
        <div class="bg-white rounded-lg border border-border shadow-sm">
          <div class="h-28 border-b border-border bg-bg-subtle flex-center">
            <div class="w-14 h-14 bg-gray-200 rounded-full animate-pulse"></div>
          </div>
          <div class="p-3">
            <div class="h-4 bg-gray-200 rounded w-3/4 animate-pulse"></div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default GridSkeleton;
