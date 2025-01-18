// import {
//   Select,
//   SelectContent,
//   SelectItem,
//   SelectTrigger,
//   SelectValue,
// } from "@/components/ui/select";
// import { sortTypes } from "@/constants";

const Sort = () => {
  // const handleSort = (value: string) => {
  //   // router.push(`${path}?sort=${value}`);
  // };

  return (
    <div>
      <p>Sort</p>
    </div>
    // <Select onValueChange={handleSort} defaultValue={sortTypes[0].value}>
    //   <SelectTrigger className="sort-select">
    //     <SelectValue placeholder={sortTypes[0].value} />
    //   </SelectTrigger>
    //   <SelectContent className="sort-select-content">
    //     {/* {sortTypes.map((sort) => (
    //       <SelectItem
    //         key={sort.label}
    //         className="shad-select-item"
    //         value={sort.value}
    //       >
    //         {sort.label}
    //       </SelectItem>
    //     ))} */}
    //   </SelectContent>
    // </Select>
  );
};

export default Sort;
