import { AppSidebar } from "@/components/home/sidebar/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { getServerSession } from "next-auth";
import { authOptions } from "./api/auth/[...nextauth]/route";
import { redirect } from "next/navigation";
import Header from "@/components/home/conversation/Header";
import Chat from "@/components/home/conversation/Chat";


export const metadata = {
  title: "Chatify",
};

export default async function Page() {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect("/auth?mode=login");
  }

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="flex flex-col h-screen">
        <Header />
        <Chat />
      </SidebarInset>
    </SidebarProvider>
  );
}
