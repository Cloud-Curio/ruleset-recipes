import React from 'react';
import { Home, Users, Video, Bookmark, Calendar, Clock, ChevronDown } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

export const LeftSidebar: React.FC = () => {
  const { user } = useAuth();

  const menuItems = [
    { icon: Home, label: 'Home', href: '/feed' },
    { icon: Users, label: 'Friends', href: '/friends', badge: '3' },
    { icon: Video, label: 'Watch', href: '/watch' },
    { icon: Bookmark, label: 'Saved', href: '/saved' },
    { icon: Calendar, label: 'Events', href: '/events' },
    { icon: Clock, label: 'Memories', href: '/memories' },
  ];

  return (
    <aside className="hidden lg:block w-80 h-[calc(100vh-56px)] sticky top-14 overflow-y-auto pb-4">
      <div className="p-2">
        {/* User Profile */}
        <Link href="/profile">
          <Button variant="ghost" className="w-full justify-start px-2 py-3 h-auto">
            <Avatar className="w-9 h-9 mr-3">
              <AvatarImage src={user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=User'} />
              <AvatarFallback>{user?.name?.[0] || 'U'}</AvatarFallback>
            </Avatar>
            <span className="font-semibold">{user?.name || 'User'}</span>
          </Button>
        </Link>

        {/* Menu Items */}
        {menuItems.map((item) => (
          <Link key={item.label} href={item.href}>
            <Button variant="ghost" className="w-full justify-start px-2 py-3 h-auto relative">
              <div className="w-9 h-9 mr-3 flex items-center justify-center bg-gray-200 dark:bg-gray-800 rounded-full">
                <item.icon className="w-5 h-5" />
              </div>
              <span className="font-medium">{item.label}</span>
              {item.badge && (
                <span className="ml-auto bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                  {item.badge}
                </span>
              )}
            </Button>
          </Link>
        ))}

        <Button variant="ghost" className="w-full justify-start px-2 py-3 h-auto">
          <div className="w-9 h-9 mr-3 flex items-center justify-center bg-gray-200 dark:bg-gray-800 rounded-full">
            <ChevronDown className="w-5 h-5" />
          </div>
          <span className="font-medium">See more</span>
        </Button>

        <Separator className="my-4" />

        {/* Shortcuts */}
        <div className="px-2 mb-2">
          <h3 className="text-gray-600 dark:text-gray-400 font-semibold text-sm px-2">
            Your Shortcuts
          </h3>
        </div>
        
        {[
          { name: 'React Developers', image: 'https://api.dicebear.com/7.x/shapes/svg?seed=React' },
          { name: 'Web Design', image: 'https://api.dicebear.com/7.x/shapes/svg?seed=Design' },
          { name: 'Tech News', image: 'https://api.dicebear.com/7.x/shapes/svg?seed=Tech' },
        ].map((shortcut) => (
          <Button key={shortcut.name} variant="ghost" className="w-full justify-start px-2 py-2 h-auto">
            <Avatar className="w-9 h-9 mr-3 rounded-lg">
              <AvatarImage src={shortcut.image} />
              <AvatarFallback>{shortcut.name[0]}</AvatarFallback>
            </Avatar>
            <span className="font-medium text-sm">{shortcut.name}</span>
          </Button>
        ))}
      </div>
    </aside>
  );
};
