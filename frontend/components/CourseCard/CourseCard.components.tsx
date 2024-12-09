"use client";

import Link from "next/link";
import Image from "next/image";
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CourseType } from "@/lib/types";
import { useState } from "react";

export default function CourseCard({ course, enrolled }: { course: CourseType, enrolled?: boolean }) {
  const [src, setSrc] = useState<string>(course.image_url ?? "/placeholder.svg")

  const courseId = course.id;

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
      <div className="relative">
        <Image
          alt="Course Image"
          onError={() => {
            setSrc("/placeholder.svg");
          }}
          className="w-full h-48 object-cover"
          height={225}
          src={src}
          style={{
            aspectRatio: "400/225",
            objectFit: "cover",
          }}
          width={400}
        />
      </div>
      <div className="p-4 flex flex-col justify-between h-max">
        <div>
          <h3 className="text-lg font-semibold mb-2">{course.title}</h3>
          <p className="text-gray-500 line-clamp-2 mb-4">{course.description}</p>
          <div className="flex gap-2 mb-4 overflow-x-auto scrollbar-none whitespace-nowrap">
            <Badge>{course.category}</Badge>
          </div>
        </div>
        <Link href={`/course/${courseId}`}>
          <Button variant={enrolled ? "outline" : "default"} className="w-full">{"View Course"}</Button>
        </Link>
      </div>
    </div>
  );
}

function StarIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
    </svg>
  )
}

function UserIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  )
}