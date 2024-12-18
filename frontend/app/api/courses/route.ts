import { formatCookies, getCookieValue } from "@/lib/api.utils";
import { jwtDecode } from "jwt-decode";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";
  const searchParams = request.nextUrl.searchParams;
  const name = searchParams.get("name") || "";
  const category = searchParams.get("category") || "";
  const available = searchParams.get("available") === "true";

  try {
    const searchApiUrl =
      process.env.NEXT_PUBLIC_SEARCH_API_URL ?? "http://search-api:8003";
    const searchUrl = `${searchApiUrl}/search?q=${encodeURIComponent(
      name
    )}&category=${encodeURIComponent(category)}&available=${available}`;

    const response = await fetch(searchUrl, {
      method: "GET",
      cache: "no-cache",
      headers: {
        Accept: "application/json",
      },
    });

    if (!response.ok) {
      throw new Error("Failed to fetch courses");
    }

    const coursesJson = await response.json();
    const courses = coursesJson?.response?.docs || [];

    return NextResponse.json({
      message: "Courses fetched successfully",
      courses,
    });
  } catch (error: any) {
    console.error("Error fetching courses:", error);
    return NextResponse.json(
      { message: error.message || "Error fetching courses", courses: [] },
      { status: 500 }
    );
  }
}

export async function POST(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_COURSES_API_URL ?? "";

  try {
    const cookieValue = getCookieValue();
    if (!cookieValue) {
      return NextResponse.json(
        { message: "Not authenticated" },
        { status: 401 }
      );
    }

    const decoded = jwtDecode(cookieValue) as any;
    if (!decoded?.admin) {
      return NextResponse.json({ message: "Not authorized" }, { status: 403 });
    }

    const data = await request.json();

    const response = await fetch(`${baseUrl}/courses`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookieValue}`,
      },
      body: JSON.stringify(data),
    });

    const responseData = await response.json();

    if (!response.ok) {
      return NextResponse.json(
        { error: responseData.error || "Failed to create course" },
        { status: response.status }
      );
    }

    return NextResponse.json(responseData);
  } catch (error: any) {
    console.error("Error creating course:", error);
    return NextResponse.json(
      { error: error.message || "Internal server error" },
      { status: 500 }
    );
  }
}
