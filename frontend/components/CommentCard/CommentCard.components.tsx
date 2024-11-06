"use client";
import { CommentType } from "@/lib/types";
import { AvatarImage, AvatarFallback, Avatar } from "@/components/ui/avatar"

export default function CommentCard({ comment }: { comment: CommentType }) {
  return (
    <div className="flex gap-4">
      <Avatar className="w-10 h-10 border">
        <AvatarImage alt="@shadcn" src="/placeholder-user.jpg" />
        <AvatarFallback>CN</AvatarFallback>
      </Avatar>
      <div className="grid gap-2">
        <div className="flex items-center gap-2">
          <h4 className="font-semibold">{comment?.user_name}</h4>
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          {comment?.comment}
        </p>
      </div>
    </div>
  );
}