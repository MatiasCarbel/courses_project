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
        <div className="absolute top-2 right-2">
          <Badge variant={course.available_seats > 0 ? "default" : "destructive"}>
            {course.available_seats > 0
              ? `${course.available_seats} seats left`
              : "No seats available"
            }
          </Badge>
        </div>
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
