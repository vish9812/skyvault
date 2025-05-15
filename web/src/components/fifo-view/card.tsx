interface CardProps {
  type: "file" | "folder";
  name: string;
}

function Card(props: CardProps) {
  return (
    <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
      <div class="p-4">
        <div class="text-lg font-semibold text-gray-900">{props.name}</div>
      </div>
    </div>
  );
}

export default Card;
