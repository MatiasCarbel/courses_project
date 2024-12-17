"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/useUser";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import { CardTitle, CardHeader, CardContent, Card } from "@/components/ui/card";
import { CourseType } from "@/lib/types";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

export default function Course({ params }: { params: { slug: string } }) {
  const [course, setCourse] = useState<CourseType>();
  const [src, setSrc] = useState<string>("/placeholder.svg");
  const [alreadyEnrolled, setAlreadyEnrolled] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [updateTitle, setUpdateTitle] = useState("");
  const [updateDescription, setUpdateDescription] = useState("");
  const [updateImage, setUpdateImage] = useState("");
  const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false);
  const router = useRouter();
  const { isAdmin } = useUser();

  useEffect(() => {
    setAlreadyEnrolled(course?.is_subscribed ?? false);
  }, [course]);

  useEffect(() => {
    updateCourses();
  }, []);

  useEffect(() => {
    if (course) {
      setUpdateTitle(course.title);
      setUpdateDescription(course.description);
      setUpdateImage(course.image_url);
    }
  }, [course]);

  const updateCourses = () => {
    fetch(`/api/courses/courseId?courseId=${params?.slug}`)
      .then((res) => res.json())
      .then((data) => {
        console.log("data: ", data);
        setCourse(data?.course);
        setSrc(data?.course?.image_url ?? "/placeholder.svg");
      });
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this course?")) return;

    setIsDeleting(true);
    try {
      const response = await fetch(`/api/courses/delete/${params?.slug}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Failed to delete course");
      }

      router.push('/home');
    } catch (error) {
      console.error("Error deleting course:", error);
      alert("Failed to delete course");
    } finally {
      setIsDeleting(false);
    }
  };

  const enroll = async () => {
    if (alreadyEnrolled || course?.available_seats === 0) return;

    const response = await fetch(`/api/courses/subscribe`, {
      method: "POST",
      body: JSON.stringify({ courseId: params?.slug }),
    }).then(async (res) => {
      const data = await res.json();
      if (res.ok) {
        router.push(`/success/${params?.slug}`);
      } else {
        alert("Error while enrolling to course: " + data?.message);
      }
    }).catch((error) => {
      console.error("Error while enrolling to course: ", error);
    });
  };

  const handleUpdate = async () => {
    setIsUpdating(true);
    try {
      const response = await fetch(`/api/courses/update/${params?.slug}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          title: updateTitle,
          description: updateDescription,
          image_url: updateImage,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to update course");
      }

      setIsUpdateModalOpen(false);
      updateCourses();
    } catch (error: any) {
      console.error("Error updating course:", error);
      alert(error.message || "Failed to update course");
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <div className="grid md:grid-cols-2 gap-6 lg:gap-12 items-start max-w-6xl px-4 mx-auto py-6 h-full">
      <div className="grid gap-4 md:gap-10 items-start">
        <Image
          alt="Course Preview"
          className="rounded-lg w-full aspect-[16/9] object-cover"
          height={450}
          onError={() => {
            setSrc("/placeholder.svg");
          }}
          src={src}
          width={800}
        />
        <div className="grid gap-4">
          <h1 className="text-3xl font-bold">{course?.title}</h1>
          <div className="grid gap-2">
            <p className="text-gray-500 dark:text-gray-400">{course?.description}</p>
            <div className="flex items-center gap-2">
              <UserIcon className="w-5 h-5 fill-muted stroke-muted-foreground" />
              <span className="text-sm text-gray-500 dark:text-gray-400">Instructor: {course?.instructor}</span>
            </div>
            <div className="flex items-center gap-2">
              <ClockIcon className="w-5 h-5 fill-muted stroke-muted-foreground" />
              <span className="text-sm text-gray-500 dark:text-gray-400">{course?.duration} hours of video</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-500 dark:text-gray-400">Available Seats: {course?.available_seats}</span>
            </div>
          </div>
        </div>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Description</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4">
            <p className="text-gray-500 dark:text-gray-400">{course?.description}</p>
          </CardContent>
        </Card>
        <div className="flex gap-4">
          <Button
            onClick={enroll}
            variant={!alreadyEnrolled ? "default" : "outline"}
            disabled={alreadyEnrolled || course?.available_seats === 0}
            size="lg"
            className="flex-1"
          >
            {!alreadyEnrolled ? "Enroll in Course" : "Already Enrolled"}
          </Button>
          {isAdmin && (
            <div className="flex gap-4">
              <Dialog open={isUpdateModalOpen} onOpenChange={setIsUpdateModalOpen}>
                <DialogTrigger asChild>
                  <Button variant="outline" size="lg">
                    Update Course
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Update Course</DialogTitle>
                  </DialogHeader>
                  <div className="grid gap-4 py-4">
                    <div className="grid gap-2">
                      <Label htmlFor="title">Title</Label>
                      <Input
                        id="title"
                        value={updateTitle}
                        onChange={(e) => setUpdateTitle(e.target.value)}
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor="description">Description</Label>
                      <Textarea
                        id="description"
                        value={updateDescription}
                        onChange={(e) => setUpdateDescription(e.target.value)}
                      />
                    </div>
                    <div className="grid gap-2">
                      <Label htmlFor="image">Image URL</Label>
                      <Input
                        id="image"
                        value={updateImage}
                        onChange={(e) => setUpdateImage(e.target.value)}
                      />
                    </div>
                  </div>
                  <Button
                    onClick={handleUpdate}
                    disabled={isUpdating}
                  >
                    {isUpdating ? "Updating..." : "Update Course"}
                  </Button>
                </DialogContent>
              </Dialog>
              <Button
                onClick={handleDelete}
                variant="destructive"
                disabled={isDeleting}
                size="lg"
              >
                {isDeleting ? "Deleting..." : "Delete Course"}
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function ClockIcon(props: React.SVGProps<SVGSVGElement>) {
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
      <circle cx="12" cy="12" r="10" />
      <polyline points="12 6 12 12 16 14" />
    </svg>
  );
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
  );
}
