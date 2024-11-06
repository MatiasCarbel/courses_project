import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const baseUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? "";

  const category = request.nextUrl.searchParams.get("category") ?? "";
  const name = request.nextUrl.searchParams.get("name") ?? "";

  const url = `${baseUrl}/courses/search?q=${name}&category=${category}`;

  const coursesReq = await fetch(url, {
    cache: "no-cache",
  });

  const coursesJson = await coursesReq.json();
  const courses = coursesJson.results;

  return NextResponse.json(
    { message: "Courses fetched", courses },
    { status: 200 }
  );
}
