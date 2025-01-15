import React from "react";
import utils, { cn } from "@/lib/utils";

interface Props {
  date: string;
  className?: string;
}

export const FormattedDateTime = ({ date, className }: Props) => {
  return (
    <p className={cn("body-1 text-light-200", className)}>
      {utils.formattedDateTime(date)}
    </p>
  );
};
export default FormattedDateTime;
