import React from "react";

interface Props {
  avatar: string | null;
  className: string;
}

const Avatar = ({ avatar, className }: Props) => {
  if (!avatar) {
    avatar = "/assets/images/avatar.png";
  }

  return (
    <img
      src={avatar}
      alt="avatar"
      width={44}
      height={44}
      className={className}
    />
  );
};

export default Avatar;
