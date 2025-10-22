import React from 'react';
import { MoreHorizontal, Search } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';

export const RightSidebar: React.FC = () => {
  const contacts = [
    { id: '1', name: 'Sarah Wilson', online: true, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah' },
    { id: '2', name: 'Michael Chen', online: true, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Michael' },
    { id: '3', name: 'Emma Davis', online: false, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma' },
    { id: '4', name: 'James Rodriguez', online: true, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=James' },
    { id: '5', name: 'Olivia Taylor', online: false, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Olivia' },
    { id: '6', name: 'William Brown', online: true, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=William' },
    { id: '7', name: 'Sophia Martinez', online: false, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sophia' },
    { id: '8', name: 'David Anderson', online: true, avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=David' },
  ];

  return (
    <aside className="hidden xl:block w-80 h-[calc(100vh-56px)] sticky top-14 overflow-y-auto pb-4">
      <div className="p-4">
        {/* Sponsored */}
        <div className="mb-4">
          <h3 className="text-gray-600 dark:text-gray-400 font-semibold mb-3">
            Sponsored
          </h3>
          <div className="space-y-3">
            <div className="flex items-start space-x-3 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800 p-2 rounded-lg">
              <img 
                src="https://api.dicebear.com/7.x/shapes/svg?seed=Ad1" 
                alt="Ad"
                className="w-24 h-24 rounded-lg object-cover"
              />
              <div className="flex-1">
                <p className="text-sm font-medium">Premium React Course</p>
                <p className="text-xs text-gray-500">learnreact.com</p>
              </div>
            </div>
          </div>
        </div>

        <Separator className="my-4" />

        {/* Contacts */}
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-gray-600 dark:text-gray-400 font-semibold">
            Contacts
          </h3>
          <div className="flex space-x-2">
            <Button variant="ghost" size="icon" className="w-8 h-8 rounded-full">
              <Search className="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="icon" className="w-8 h-8 rounded-full">
              <MoreHorizontal className="w-4 h-4" />
            </Button>
          </div>
        </div>

        <div className="space-y-1">
          {contacts.map((contact) => (
            <Button 
              key={contact.id} 
              variant="ghost" 
              className="w-full justify-start px-2 py-2 h-auto relative hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              <div className="relative">
                <Avatar className="w-9 h-9 mr-3">
                  <AvatarImage src={contact.avatar} />
                  <AvatarFallback>{contact.name[0]}</AvatarFallback>
                </Avatar>
                {contact.online && (
                  <span className="absolute bottom-0 right-2 w-3 h-3 bg-green-500 border-2 border-white dark:border-gray-900 rounded-full"></span>
                )}
              </div>
              <span className="font-medium text-sm">{contact.name}</span>
            </Button>
          ))}
        </div>
      </div>
    </aside>
  );
};
