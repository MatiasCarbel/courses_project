"use client";
import { useEffect, useState } from "react"
import { CourseType } from "@/lib/types";
import CourseCard from "@/components/CourseCard/CourseCard.components";

export default function Component() {
  const [courses, setCourses] = useState<CourseType[]>([])

  useEffect(() => {
    fetch(`/api/courses/myCourses`)
      .then((res) => res.json())
      .then((data) => {
        console.log("data: ", data);
        setCourses(data?.data?.courses ?? [])
      })
  }, [])

  return (
    <main className="container mx-auto py-8 px-4 md:px-6 lg:px-8">
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">My Courses</h1>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {Array.isArray(courses) && courses.map((course) => (
          <CourseCard enrolled={true} key={course.id} course={course} />
        ))}
        {(!courses || courses.length === 0) && (
          <div className="col-span-full">
            <p className="text-center text-gray-500">No courses found</p>
          </div>
        )}
      </div>
    </main>
  )
}
