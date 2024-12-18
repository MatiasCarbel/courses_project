"use client";
import { CourseType } from "@/lib/types";
import Image from "next/image";
import Link from "next/link"
import { useEffect, useState } from "react";

export default function Component({ params }: { params: { slug: string } }) {
  const [course, setCourse] = useState<CourseType>()
  const [src, setSrc] = useState<string>("/placeholder.svg")

  useEffect(() => {
    fetch(`/api/courses/courseId?courseId=${params?.slug}`)
      .then((res) => res.json())
      .then((data) => {
        setCourse(data?.course?.data);
        setSrc(data?.course?.data?.image_url ?? "/placeholder.svg");
      })
  }, [params?.slug])

  return (
    <div className="flex flex-col">
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32">
          <div className="container grid items-center gap-6 px-4 md:px-6 lg:grid-cols-2 lg:gap-10">
            <Image
              src={src}
              width="600"
              height="500"
              alt="Course Image"
              className="mx-auto aspect-[3/2] overflow-hidden rounded-xl object-cover"
            />
            <div className="space-y-4">
              <div className="inline-block rounded-lg bg-muted px-3 py-1 text-sm">Congratulations!</div>
              <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
                You&#39;ve enrolled in {course?.title ?? "Course Name"}
              </h1>
              <p className="max-w-[600px] text-muted-foreground md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed">
                We&#39;re excited to have you join our program.
              </p>
              <Link
                href="/home"
                className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-8 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50"
                prefetch={false}
              >
                Back to Home
              </Link>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}