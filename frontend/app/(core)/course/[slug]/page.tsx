"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/useUser";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import CommentCard from "@/components/CommentCard/CommentCard.components";
import { CardTitle, CardHeader, CardContent, Card } from "@/components/ui/card";
import { CommentType, CourseType } from "@/lib/types";

export default function Course({ params }: { params: { slug: string } }) {
  const [course, setCourse] = useState<CourseType>();
  const [comments, setComments] = useState<CommentType[]>([]);
  const [src, setSrc] = useState<string>("/placeholder.svg");
  const [alreadyEnrolled, setAlreadyEnrolled] = useState(false);
  const { isAdmin } = useUser();
  const router = useRouter();

  useEffect(() => {
    setAlreadyEnrolled(course?.is_subscribed ?? false);
  }, [course]);

  useEffect(() => {
    updateCourses();
  }, []);

  const updateCourses = () => {
    fetch(`/api/courses/courseId?courseId=${params?.slug}`)
      .then((res) => res.json())
      .then((data) => {
        setCourse(data?.course);
        setComments(data?.comments.slice(0, 4));
        setSrc(data?.course?.ImageURL ?? "/placeholder.svg");
      });
  };

  const uploadResource = async () => {
    const fileInput = document.createElement("input");
    fileInput.type = "file";
    fileInput.accept = "application/pdf";
    fileInput.click();

    fileInput.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (!file) return;

      const formData = new FormData();
      formData.append("file", file);
      formData.append("courseId", params?.slug);

      const response = await fetch(`/api/courses/resources`, {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        alert("Error while uploading resource.");
        return;
      }
      const responseText = await response.text();

      let responseJson;
      try {
        responseJson = JSON.parse(responseText);
      } catch (error) {
        alert("Error while parsing server response.");
        return;
      }

      if (response.ok) {
        alert("Resource uploaded successfully.");
      } else {
        alert("Error while uploading resource.");
      }
    };
  };

  const downloadResources = async () => {
    const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";
    const url = `${baseUrl}/download/${params?.slug}`;
    console.log(url);

    window.location.href = url;
    // const resourceReq = await fetch(url);
    // const resourceJson = await resourceReq.json();

    // console.log(resourceJson);

    // if (!resourceReq.ok) {
    //   alert("Error while downloading resources.");
    //   return;
    // }
  };

  const addComment = () => {
    const comment = prompt("Enter your comment:");

    if (comment) {
      fetch(`/api/courses/comment`, {
        method: "POST",
        body: JSON.stringify({ courseId: params?.slug, comment }),
      })
        .then((res) => res.json())
        .then((res) => {
          updateCourses();
        });
    }
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString();
  };

  const enroll = async () => {
    if (alreadyEnrolled) return;

    const response = await fetch(`/api/courses/subscribe`, {
      method: "POST",
      body: JSON.stringify({ courseId: params?.slug }),
    });

    if (response.ok) {
      const newCourse = { ...course, is_subscribed: true } as CourseType;
      setCourse(newCourse);
      router.push(`/success/${params?.slug}`);
    } else {
      alert("Error while enrolling to course.");
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
          <h1 className="text-3xl font-bold">{course?.course_name}</h1>
          <div className="grid gap-2">
            <p className="text-gray-500 dark:text-gray-400">{course?.description}</p>
            <div className="flex items-center gap-2">
              <UserIcon className="w-5 h-5 fill-muted stroke-muted-foreground" />
              <span className="text-sm text-gray-500 dark:text-gray-400">Instructor: {course?.instructor_name}</span>
            </div>
            <div className="flex items-center gap-2">
              <ClockIcon className="w-5 h-5 fill-muted stroke-muted-foreground" />
              <span className="text-sm text-gray-500 dark:text-gray-400">{course?.length} hours of video</span>
            </div>
            <div className="flex items-center gap-2">
              <CalendarIcon className="w-5 h-5 fill-muted stroke-muted-foreground" />
              <span className="text-sm text-gray-500 dark:text-gray-400">Last updated: {formatDate(course?.LastUpdated ?? "")}</span>
            </div>
          </div>
        </div>
      </div>
      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Comments</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4">
            {comments?.map((comment) => (
              <CommentCard key={comment.comment_id} comment={comment} />
            ))}
            <Button onClick={addComment} variant={"secondary"} size="lg">Add Comment</Button>
          </CardContent>
        </Card>
        <Button onClick={enroll} variant={!alreadyEnrolled ? "default" : "outline"} disabled={alreadyEnrolled} size="lg">{!alreadyEnrolled ? "Enroll in Course" : "Already Enrolled"}</Button>
        {alreadyEnrolled && (
          <Button onClick={downloadResources} variant={"default"} size="lg">Download Course Resources</Button>
        )}
        {isAdmin && (
          <Button onClick={uploadResource} variant={"default"} size="lg">Upload Course Resource</Button>
        )}
      </div>
    </div>
  );
}

function CalendarIcon(props: React.SVGProps<SVGSVGElement>) {
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
      <path d="M8 2v4" />
      <path d="M16 2v4" />
      <rect width="18" height="18" x="3" y="4" rx="2" />
      <path d="M3 10h18" />
    </svg>
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
