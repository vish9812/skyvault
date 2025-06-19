function ListSkeleton() {
  return (
    <div class="border border-border rounded-lg overflow-hidden max-w-screen-xl mx-auto">
      {/* Header */}
      <div class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-3 bg-bg-muted border-b border-border font-medium text-sm text-neutral-light">
        <div></div>
        <div>Name</div>
        <div class="text-right">Size</div>
        <div class="text-right">Last modified</div>
        <div></div>
      </div>

      {/* Skeleton rows */}
      {Array.from({ length: 7 }).map(() => (
        <div class="grid grid-cols-[2rem_1fr_6rem] md:grid-cols-[2rem_1fr_6rem_9rem] items-center py-3 px-4 border-t border-border first:border-t-0">
          <div class="h-6 w-6 bg-gray-200 rounded-full animate-pulse"></div>
          <div class="min-w-0">
            <div class="h-5 bg-gray-200 rounded-sm w-1/3 animate-pulse"></div>
          </div>
          <div class="text-right">
            <div class="h-4 bg-gray-200 rounded-sm w-12 ml-auto animate-pulse"></div>
          </div>
          <div class="text-right">
            <div class="h-4 bg-gray-200 rounded-sm w-24 ml-auto animate-pulse"></div>
          </div>
          <div></div>
        </div>
      ))}
    </div>
  );
}

export default ListSkeleton;
