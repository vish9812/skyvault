function ListSkeleton() {
  return (
    <div class="flex flex-col border border-border rounded-lg overflow-hidden">
      {/* Header */}
      <div class="flex items-center py-3 px-4 bg-bg-muted border-b border-border font-medium text-sm text-neutral-light">
        <div class="w-6 mr-3"></div>
        <div class="flex-1 min-w-0">Name</div>
        <div class="w-24 text-right">Size</div>
        <div class="w-36 text-right">Last modified</div>
        <div class="w-8"></div>
      </div>

      {/* Skeleton rows */}
      {Array.from({ length: 7 }).map(() => (
        <div class="flex items-center py-3 px-4 border-t border-border first:border-t-0">
          <div class="w-6 h-6 bg-gray-200 rounded-full mr-3 animate-pulse"></div>
          <div class="flex-1 min-w-0">
            <div class="h-5 bg-gray-200 rounded w-1/3 animate-pulse"></div>
          </div>
          <div class="w-24 text-right">
            <div class="h-4 bg-gray-200 rounded w-12 ml-auto animate-pulse"></div>
          </div>
          <div class="w-36 text-right">
            <div class="h-4 bg-gray-200 rounded w-24 ml-auto animate-pulse"></div>
          </div>
          <div class="w-8"></div>
        </div>
      ))}
    </div>
  );
}

export default ListSkeleton;
