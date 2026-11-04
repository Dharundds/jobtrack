import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  Columns3,
  List,
  Settings,
  Radio,
} from "lucide-react";
import { cn } from "@/lib/utils";

const links = [
  { to: "/", label: "Overview", icon: LayoutDashboard, end: true },
  { to: "/board", label: "Board", icon: Columns3 },
  { to: "/applications", label: "Rows", icon: List },
  { to: "/settings", label: "Settings", icon: Settings },
] as const;

function linkEnd(link: (typeof links)[number]) {
  return "end" in link && link.end === true;
}

export function Sidebar() {
  return (
    <aside
      className="hidden w-52 shrink-0 flex-col bg-rail text-rail-foreground md:flex"
      aria-label="Main navigation"
    >
      <div className="flex h-14 items-center gap-2.5 border-b border-white/10 px-4">
        <Radio className="h-5 w-5 text-accent" aria-hidden />
        <div>
          <span className="block text-sm font-semibold tracking-wide">jobtrack</span>
          <span className="font-data text-[10px] uppercase tracking-widest text-rail-muted">
            local · read-only
          </span>
        </div>
      </div>
      <nav className="flex flex-1 flex-col gap-0.5 p-2">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            end={linkEnd(link)}
            className={({ isActive }) =>
              cn(
                "relative flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium transition-colors",
                isActive
                  ? "bg-white/10 text-rail-foreground"
                  : "text-rail-muted hover:bg-white/5 hover:text-rail-foreground",
              )
            }
          >
            {({ isActive }) => (
              <>
                {isActive ? (
                  <span
                    className="absolute bottom-2 left-0 top-2 w-0.5 rounded-full bg-accent"
                    aria-hidden
                  />
                ) : null}
                <link.icon className="h-4 w-4 shrink-0" aria-hidden />
                {link.label}
              </>
            )}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}

export function MobileNav() {
  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-40 flex border-t border-border bg-surface-elevated pb-[env(safe-area-inset-bottom)] shadow-board md:hidden"
      aria-label="Mobile navigation"
    >
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          end={linkEnd(link)}
          className={({ isActive }) =>
            cn(
              "flex flex-1 flex-col items-center gap-0.5 py-2 text-[10px] font-semibold uppercase tracking-wide",
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
