import { formatCookies } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";
  const coursesApiUrl =
    process.env.NEXT_PUBLIC_COURSES_API_URL ?? "http://courses-api:8002";

  // Get courseId from query params
  const courseId = request.nextUrl.searchParams.get("courseId");

  try {
    // Fetch course details from Search API
    console.log("courseId: ", courseId);
    const searchResponse = await fetch(`${coursesApiUrl}/courses/${courseId}`, {
      headers: {
        Accept: "application/json",
      },
      cache: "no-cache",
    });

    const course = await searchResponse.json();
    console.log("course: ", course);

    if (!course) {
      return NextResponse.json(
        { message: "Course not found" },
        { status: 404 }
      );
    }

    // Fetch enrollment status
    // const cookie = formatCookies();
    // const enrollmentResponse = await fetch(
    //   `${coursesApiUrl}/enrollments/check/${courseId}`,
    //   {
    //     headers: {
    //       cookie: cookie,
    //     },
    //   }
    // );

    // const enrollmentData = await enrollmentResponse.json();
    // console.log("enrollmentData: ", enrollmentData);

    // Fetch comments
    // const commentsResponse = await fetch(`${baseUrl}/comments/${courseId}`, {
    //   headers: {
    //     Accept: "application/json",
    //   },
    // });
    // const comments = await commentsResponse.json();
    // console.log("comments: ", comments);

    return NextResponse.json(
      {
        message: "Course details fetched",
        course: {
          ...course,
          // is_subscribed:
          //   enrollmentData?.message === "User is enrolled in this course",
        },
        // comments: comments,
      },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Error fetching course details:", error);
    return NextResponse.json(
      { message: "Error fetching course details", error: error.message },
      { status: 500 }
    );
  }
}
