import { formatCookies } from "@/lib/api.utils";
import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";
import { jwtDecode } from "jwt-decode";

export async function GET(request: NextRequest) {
  const cookiesStore = cookies();
  const authCookie = cookiesStore.get("auth");

  const cookieValue = authCookie?.value;
  if (!cookieValue)
    return NextResponse.json(
      { message: "Not authenticated.", shouldLogin: true },
      { status: 200 }
    );

  const decoded = jwtDecode(cookieValue ?? "") as any;
  const userId = decoded?.user_id;

  const baseUrl = process.env.NEXT_PUBLIC_USERS_API_URL ?? "";
  const searchApiUrl =
    process.env.NEXT_PUBLIC_SEARCH_API_URL ?? "http://search-api:8003";

  // Get courseId from query params
  const courseId = request.nextUrl.searchParams.get("courseId");

  try {
    // Fetch course details from Search API
    const searchResponse = await fetch(
      `${searchApiUrl}/search?q=id:"${courseId}"&wt=json`,
      {
        headers: {
          Accept: "application/json",
        },
        cache: "no-cache",
      }
    );

    if (!searchResponse.ok) {
      throw new Error(
        `Search API responded with status: ${searchResponse.status}`
      );
    }

    const searchData = await searchResponse.json();
    const course = searchData?.response?.docs?.[0];

    if (!course) {
      return NextResponse.json(
        { message: "Course not found" },
        { status: 404 }
      );
    }

    // Fetch enrollment status
    const cookie = formatCookies();
    const enrollmentResponse = await fetch(
      `${baseUrl}/enrollments/check/${courseId}`,
      {
        headers: {
          cookie: cookie,
        },
      }
    );

    const enrollmentData = await enrollmentResponse.json();

    // Fetch comments
    const commentsResponse = await fetch(`${baseUrl}/comments/${courseId}`, {
      headers: {
        cookie: cookie,
      },
    });
    const comments = await commentsResponse.json();

    return NextResponse.json(
      {
        message: "Course details fetched",
        course: {
          ...course,
          is_subscribed:
            enrollmentData?.message === "User is enrolled in this course",
        },
        comments: comments,
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
