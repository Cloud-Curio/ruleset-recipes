import React from 'react';
import { Plus } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Card } from '@/components/ui/card';
import { useAuth } from '@/contexts/AuthContext';

export const Stories: React.FC = () => {
  const { user } = useAuth();

  const stories = [
    { id: '1', name: 'Sarah Wilson', image: 'https://picsum.photos/seed/sarah/400/600', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah', hasNew: true },
    { id: '2', name: 'Michael Chen', image: 'https://picsum.photos/seed/michael/400/600', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Michael', hasNew: true },
    { id: '3', name: 'Emma Davis', image: 'https://picsum.photos/seed/emma/400/600', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma', hasNew: true },
    { id: '4', name: 'James Rodriguez', image: 'https://picsum.photos/seed/james/400/600', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=James', hasNew: false },
  ];

  return (
    <div className="mb-4">
      <div className="flex space-x-2 overflow-x-auto pb-2 scrollbar-hide">
        {/* Create Story Card */}
        <Card className="relative flex-shrink-0 w-28 h-48 overflow-hidden cursor-pointer group hover:shadow-lg transition-shadow">
          <img 
            src={user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=User'} 
            alt="Create story"
            className="w-full h-32 object-cover"
          />
          <div className="absolute top-20 left-1/2 transform -translate-x-1/2 w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center border-4 border-white dark:border-gray-900">
            <Plus className="w-5 h-5 text-white" />
          </div>
          <div className="absolute bottom-0 left-0 right-0 p-2 text-center">
            <p className="text-xs font-semibold">Create Story</p>
          </div>
        </Card>

        {/* Story Cards */}
        {stories.map((story) => (
          <Card 
            key={story.id} 
            className="relative flex-shrink-0 w-28 h-48 overflow-hidden cursor-pointer group hover:shadow-lg transition-shadow"
          >
            <img 
              src={story.image} 
              alt={story.name}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
            />
            <div className="absolute inset-0 bg-gradient-to-b from-black/50 via-transparent to-black/70"></div>
            <div className={`absolute top-3 left-3 w-10 h-10 rounded-full ${story.hasNew ? 'ring-4 ring-blue-600' : 'ring-2 ring-gray-400'}`}>
              <Avatar className="w-full h-full">
                <AvatarImage src={story.avatar} />
                <AvatarFallback>{story.name[0]}</AvatarFallback>
              </Avatar>
            </div>
            <div className="absolute bottom-0 left-0 right-0 p-2">
              <p className="text-xs font-semibold text-white drop-shadow-lg">
                {story.name}
              </p>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};
