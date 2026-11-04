import { useState } from "react";
import { Outlet } from "react-router-dom";
import { MobileNav, Sidebar } from "@/components/layout/Nav";
import { TopBar } from "@/components/layout/TopBar";

export function AppLayout() {
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <div className="flex min-h-screen bg-surface">
      <Sidebar />
      <div className="flex min-h-screen min-w-0 flex-1 flex-col pb-16 md:pb-0">
        <TopBar searchQuery={searchQuery} onSearchChange={setSearchQuery} />
        <main className="flex-1 px-4 py-6 md:px-8 md:py-8">
          <Outlet context={{ searchQuery }} />
        </main>
      </div>
      <MobileNav />
    </div>
  );
}
