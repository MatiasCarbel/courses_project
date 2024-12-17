"use client";
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { DropdownMenuTrigger, DropdownMenuContent, DropdownMenu, DropdownMenuRadioGroup, DropdownMenuRadioItem } from "@/components/ui/dropdown-menu"
import { useEffect, useState } from "react"
import { CourseType } from "@/lib/types";
import CourseCard from "@/components/CourseCard/CourseCard.components";
import { useDebounce } from "@/hooks/useDebounce";
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"

export default function Component() {
  const [courses, setCourses] = useState<CourseType[]>([])
  const [category, setCategory] = useState<string>("")
  const [search, setSearch] = useState<string>("")
  const [isLoading, setIsLoading] = useState(false)
  const [availableOnly, setAvailableOnly] = useState(false)
  const debouncedSearch = useDebounce(search, 500)

  const fetchCourses = async (searchTerm: string, categoryFilter: string) => {
    setIsLoading(true)
    try {
      const res = await fetch(`/api/courses?name=${searchTerm}&category=${categoryFilter}&available=${availableOnly}`)
      const data = await res.json()
      if (data?.courses) {
        // Get course IDs
        const courseIds = data.courses.map((course: CourseType) => course.id)

        // Fetch availabilities
        const availRes = await fetch('/api/courses/availability', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(courseIds)
        })

        const availData = await availRes.json()

        // Merge availability data with courses
        const coursesWithAvailability = data.courses.map((course: CourseType) => ({
          ...course,
          available_seats: availData[course.id as keyof typeof availData] || 0
        }))

        setCourses(coursesWithAvailability)
      }
    } catch (error) {
      console.error('Error fetching courses:', error)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchCourses(debouncedSearch, category)
  }, [debouncedSearch, category, availableOnly])

  return (
    <main className="container mx-auto py-8 px-4 md:px-6 lg:px-8">
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">All Courses</h1>
        <div className="flex items-center gap-4 mt-4 md:mt-0">
          <div className="flex items-center space-x-2">
            <Checkbox
              id="available"
              checked={availableOnly}
              onCheckedChange={(checked: boolean) => setAvailableOnly(checked)}
            />
            <Label htmlFor="available">Available Only</Label>
          </div>
          <Input
            className="w-full md:w-auto"
            value={search}
            onChange={e => setSearch(e.currentTarget.value)}
            placeholder="Search courses..."
            type="text"
          />
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline">
                {category || "All Categories"}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuRadioGroup value={category} onValueChange={setCategory}>
                <DropdownMenuRadioItem value="">All</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="web-development">Web Development</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="mobile-development">Mobile Development</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="data-science">Data Science</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="design">Design</DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="business">Business</DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-8">Loading...</div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {courses.map((course) => (
            <CourseCard key={course.id} course={course} />
          ))}
          {courses.length === 0 && !isLoading && (
            <div className="col-span-full text-center py-8 text-gray-500">
              No courses found
            </div>
          )}
        </div>
      )}
    </main>
  )
}
