function CardsSkeleton() {
  return (
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {Array.from({ length: 10 }).map(() => (
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
          <div class="flex flex-col items-center">
            <div class="w-10 h-10 bg-gray-200 rounded mb-3 animate-pulse"></div>
            <div class="h-4 bg-gray-200 rounded w-20 animate-pulse"></div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default CardsSkeleton;
