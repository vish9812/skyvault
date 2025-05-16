import CreateUpload from "@sv/components/createUpload";

function EmptyState() {
  return (
    <div class="flex-center flex-col text-center py-16 bg-white rounded-lg border border-gray-200 shadow-sm">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-16 h-16 text-primary"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M2.25 15a4.5 4.5 0 0 0 4.5 4.5H18a3.75 3.75 0 0 0 1.332-7.257 3 3 0 0 0-3.758-3.848 5.25 5.25 0 0 0-10.233 2.33A4.502 4.502 0 0 0 2.25 15Z"
        />
      </svg>
      <h3 class="text-lg font-medium text-gray-900">No files or folders yet</h3>
      <p class="mt-2 text-sm text-gray-500 max-w-md mx-auto">
        Upload files or create folders to get started with your secure cloud
        storage.
      </p>
      <div class="mt-6">
        <CreateUpload />
      </div>
    </div>
  );
}

export default EmptyState;
