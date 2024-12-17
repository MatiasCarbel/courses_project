"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { useUser } from "@/hooks/useUser";

export default function Component() {
  const [courseImage, setCourseImage] = useState<string>("");
  const [courseName, setCourseName] = useState<string>("");
  const [courseDescription, setCourseDescription] = useState<string>("");
  const [courseDuration, setCourseDuration] = useState<number>(0);
  const [courseCategory, setCourseCategory] = useState<string>("");
  const [courseInstructor, setCourseInstructor] = useState<string>("");
  const [availableSeats, setAvailableSeats] = useState<number>(0);

  const [isValid, setIsValid] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  const router = useRouter();
  const { isAdmin, isLoading } = useUser();

  // Redirect if user is not admin
  useEffect(() => {
    if (!isAdmin && !isLoading) {
      router.push('/home');
    }
  }, [isAdmin, router, isLoading]);

  useEffect(() => {
    setIsValid(
      courseImage.length > 0 &&
      courseName.length > 0 &&
      courseDescription.length > 0 &&
      courseDuration > 0 &&
      courseCategory.length > 0 &&
      courseInstructor.length > 0 &&
      availableSeats > 0
    );
  }, [courseImage, courseName, courseDescription, courseDuration, courseCategory, courseInstructor, availableSeats]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (isSubmitting || !isValid) return;
    setIsSubmitting(true);

    try {
      const response = await fetch("/api/courses/createCourse", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          courseImage,
          courseName,
          courseDescription,
          courseDuration,
          courseCategory,
          courseInstructor,
          availableSeats,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to create course');
      }

      router.push('/home');
    } catch (error) {
      console.error('Error creating course:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isAdmin) {
    return null;
  }

  return (
    <Card className="w-full max-w-2xl m-auto">
      <CardHeader>
        <CardTitle className="text-[#9f33e9]">Create a New Course</CardTitle>
        <CardDescription>Fill out the form to add a new course.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="grid gap-6 w-full">
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="image">Course Image URL</Label>
              <Input
                id="image"
                placeholder="Enter image URL"
                value={courseImage}
                onChange={(e) => setCourseImage(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Course Name</Label>
              <Input
                id="name"
                placeholder="Enter course name"
                value={courseName}
                onChange={(e) => setCourseName(e.target.value)}
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Course Description</Label>
            <Textarea
              id="description"
              placeholder="Enter course description"
              value={courseDescription}
              onChange={(e) => setCourseDescription(e.target.value)}
            />
          </div>
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="instructor">Course Instructor</Label>
              <Input
                id="instructor"
                placeholder="Enter instructor name"
                value={courseInstructor}
                onChange={(e) => setCourseInstructor(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="category">Course Category</Label>
              <Select
                value={courseCategory}
                onValueChange={(value) => setCourseCategory(value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="web-development">Web Development</SelectItem>
                  <SelectItem value="mobile-development">Mobile Development</SelectItem>
                  <SelectItem value="data-science">Data Science</SelectItem>
                  <SelectItem value="design">Design</SelectItem>
                  <SelectItem value="business">Business</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="duration">Course Duration (hours)</Label>
              <Input
                id="duration"
                type="number"
                min="1"
                placeholder="Enter duration"
                value={courseDuration}
                onChange={(e) => setCourseDuration(Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="seats">Available Seats</Label>
              <Input
                id="seats"
                type="number"
                min="1"
                placeholder="Enter available seats"
                value={availableSeats}
                onChange={(e) => setAvailableSeats(Number(e.target.value))}
              />
            </div>
          </div>
          <Button type="submit" className="justify-self-end" disabled={!isValid || isSubmitting}>
            Create Course
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}