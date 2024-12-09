import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const id = searchParams.get("id");

  // If ID is provided, fetch specific course
  if (id) {
    const baseUrl =
      process.env.NEXT_PUBLIC_SEARCH_API_URL ?? "http://search-api:8003";
    try {
      const courseUrl = `${baseUrl}/course/${id}`;
      const courseReq = await fetch(courseUrl, {
        method: "GET",
        cache: "no-cache",
        headers: {
          Accept: "application/json",
        },
      });

      if (!courseReq.ok) {
        throw new Error(`Failed to fetch course: ${courseReq.statusText}`);
      }

      const course = await courseReq.json();
      return NextResponse.json(
        { message: "Course fetched", course },
        { status: 200 }
      );
    } catch (error: any) {
      console.error("Course API error:", error);
      return NextResponse.json(
        { message: "Error fetching course", error: error.message },
        { status: 500 }
      );
    }
  }

  // Existing search logic
  const name = searchParams.get("name") || "";
  const category = searchParams.get("category") || "";

  const baseUrl =
    process.env.NEXT_PUBLIC_SEARCH_API_URL ?? "http://search-api:8003";

  try {
    const searchUrl = `${baseUrl}/search?q=${encodeURIComponent(
      name
    )}&category=${encodeURIComponent(category)}`;

    console.log("searchUrl: ", searchUrl);

    const coursesReq = await fetch(searchUrl, {
      method: "GET",
      cache: "no-cache",
      headers: {
        Accept: "application/json",
      },
    });

    if (!coursesReq.ok) {
      throw new Error(`Failed to fetch courses: ${coursesReq.statusText}`);
    }

    const coursesJson = await coursesReq.json();
    const courses = coursesJson?.response?.docs || [];

    return NextResponse.json(
      { message: "Courses fetched", courses },
      { status: 200 }
    );
  } catch (error: any) {
    console.error("Courses API error:", error);
    return NextResponse.json(
      { message: "Error fetching courses", error: error.message },
      { status: 500 }
    );
  }
}
