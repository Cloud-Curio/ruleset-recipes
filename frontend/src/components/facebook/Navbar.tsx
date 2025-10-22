import React from 'react';
import Link from 'next/link';
import { 
  Home, 
  Users, 
  Video, 
  Store, 
  Bell, 
  MessageCircle, 
  Search,
  Menu,
  Plus
} from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/contexts/AuthContext';

export const Navbar: React.FC = () => {
  const { user } = useAuth();

  return (
    <nav className="sticky top-0 z-50 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 shadow-sm">
      <div className="max-w-full px-4">
        <div className="flex items-center justify-between h-14">
          {/* Left Section */}
          <div className="flex items-center space-x-2">
            <Link href="/feed" className="flex items-center space-x-2">
              <div className="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center">
                <span className="text-white font-bold text-xl">f</span>
              </div>
            </Link>
            <div className="hidden md:flex relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
              <Input 
                type="text" 
                placeholder="Search Facebook" 
                className="pl-10 w-60 bg-gray-100 dark:bg-gray-800 border-none"
              />
            </div>
          </div>

          {/* Center Section - Navigation */}
          <div className="hidden md:flex items-center space-x-2">
            <Link href="/feed">
              <Button variant="ghost" size="icon" className="relative w-24 h-12 rounded-lg">
                <Home className="w-6 h-6" />
                <div className="absolute bottom-0 left-0 right-0 h-1 bg-blue-600 rounded-t-lg"></div>
              </Button>
            </Link>
            <Link href="/friends">
              <Button variant="ghost" size="icon" className="w-24 h-12 rounded-lg">
                <Users className="w-6 h-6" />
              </Button>
            </Link>
            <Link href="/videos">
              <Button variant="ghost" size="icon" className="w-24 h-12 rounded-lg">
                <Video className="w-6 h-6" />
              </Button>
            </Link>
            <Link href="/marketplace">
              <Button variant="ghost" size="icon" className="w-24 h-12 rounded-lg">
                <Store className="w-6 h-6" />
              </Button>
            </Link>
          </div>

          {/* Right Section */}
          <div className="flex items-center space-x-2">
            <Button variant="ghost" size="icon" className="md:flex hidden rounded-full">
              <Menu className="w-5 h-5" />
            </Button>
            <Button variant="ghost" size="icon" className="md:flex hidden rounded-full">
              <MessageCircle className="w-5 h-5" />
            </Button>
            <Button variant="ghost" size="icon" className="relative rounded-full">
              <Bell className="w-5 h-5" />
              <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
            </Button>
            <Link href="/profile">
              <Avatar className="w-9 h-9 cursor-pointer">
                <AvatarImage src={user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=User'} />
                <AvatarFallback>{user?.name?.[0] || 'U'}</AvatarFallback>
              </Avatar>
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
};
