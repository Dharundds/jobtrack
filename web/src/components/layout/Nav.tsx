import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  Columns3,
  List,
  Settings,
  Briefcase,
} from "lucide-react";
import { cn } from "@/lib/utils";

const links = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard, end: true },
  { to: "/board", label: "Board", icon: Columns3 },
  { to: "/applications", label: "List", icon: List },
  { to: "/settings", label: "Settings", icon: Settings },
] as const;

function linkEnd(link: (typeof links)[number]) {
  return "end" in link && link.end === true;
}

export function Sidebar() {
  return (
    <aside
      className="hidden w-56 shrink-0 flex-col border-r border-border bg-surface-muted md:flex"
      aria-label="Main navigation"
    >
      <div className="flex h-14 items-center gap-2 border-b border-border px-4">
        <Briefcase className="h-5 w-5 text-accent" aria-hidden />
        <span className="font-semibold tracking-tight">Jobtrack</span>
      </div>
      <nav className="flex flex-1 flex-col gap-1 p-3">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            end={linkEnd(link)}
            className={({ isActive }) =>
              cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-surface-elevated text-foreground shadow-sm"
                  : "text-foreground-muted hover:bg-surface-elevated/60 hover:text-foreground",
              )
            }
          >
            <link.icon className="h-4 w-4 shrink-0" aria-hidden />
            {link.label}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}

export function MobileNav() {
  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-40 flex border-t border-border bg-surface-elevated pb-[env(safe-area-inset-bottom)] md:hidden"
      aria-label="Mobile navigation"
    >
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          end={linkEnd(link)}
          className={({ isActive }) =>
            cn(
              "flex flex-1 flex-col items-center gap-0.5 py-2 text-[10px] font-medium",
              isActive ? "text-accent" : "text-foreground-muted",
            )
          }
        >
          <link.icon className="h-5 w-5" aria-hidden />
          {link.label}
        </NavLink>
      ))}
    </nav>
  );
}
