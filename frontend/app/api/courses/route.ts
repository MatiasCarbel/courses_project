import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const name = searchParams.get("name") || "";
  const category = searchParams.get("category") || "";

  const baseUrl =
    process.env.NEXT_PUBLIC_BASE_API_URL ?? "http://search-api:8003";

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

    console.log("coursesReq: ", coursesReq);

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
