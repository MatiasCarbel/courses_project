import routes from "@/lib/routes";
import { UserType } from "@/lib/types";
import { useRouter } from "next/navigation";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";

export function useUser() {
  const [user, setUser] = useState<UserType>();
  const [isAdmin, setIsAdmin] = useState(false);
  const [isAuthed, setIsAuthed] = useState(false);
  const router = useRouter();
  const pathname = usePathname();

  const refreshUser = async () => {
    const response = await fetch("/api/auth/user");
    const userJson = await response.json();
    console.log(userJson);

    if (userJson?.shouldLogin) {
      setUser(undefined);
      return;
    }

    setUser(userJson?.user);
    setIsAuthed(true);
    setIsAdmin(userJson?.user?.usertype);
  };

  const logout = () => {
    console.log("logging out");
    fetch("/api/auth/logout").then((response) => {
      if (response.ok) {
        setUser(undefined);
        setIsAuthed(false);
        setIsAdmin(false);
        router.push(routes.login);
      }
    });
  };

  // on route change refresh user.
  useEffect(() => {
    refreshUser();
  }, [pathname]);

  return { user, isAdmin, isAuthed, refreshUser, logout };
}
