import React, { useState } from 'react';
import { Image, Video, Smile, MapPin } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { useAuth } from '@/contexts/AuthContext';

export const PostComposer: React.FC = () => {
  const { user } = useAuth();
  const [postText, setPostText] = useState('');

  const handlePost = () => {
    if (postText.trim()) {
      // Handle post creation
      console.log('Creating post:', postText);
      setPostText('');
    }
  };

  return (
    <Card className="mb-4">
      <CardContent className="p-4">
        <div className="flex space-x-3 mb-3">
          <Avatar className="w-10 h-10">
            <AvatarImage src={user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=User'} />
            <AvatarFallback>{user?.name?.[0] || 'U'}</AvatarFallback>
          </Avatar>
          <div 
            className="flex-1 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-full px-4 py-3 cursor-pointer transition-colors"
            onClick={() => document.getElementById('post-textarea')?.focus()}
          >
            <p className="text-gray-500 dark:text-gray-400">
              What's on your mind, {user?.name?.split(' ')[0] || 'User'}?
            </p>
          </div>
        </div>

        {postText && (
          <div className="mb-3">
            <textarea
              id="post-textarea"
              value={postText}
              onChange={(e) => setPostText(e.target.value)}
              placeholder={`What's on your mind, ${user?.name?.split(' ')[0] || 'User'}?`}
              className="w-full p-3 bg-transparent border-none outline-none resize-none text-lg"
              rows={3}
              autoFocus
            />
          </div>
        )}

        <Separator className="my-3" />

        <div className="flex justify-between items-center">
          <div className="flex space-x-1 md:space-x-2">
            <Button variant="ghost" size="sm" className="text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800">
              <Video className="w-5 h-5 mr-1 md:mr-2 text-red-500" />
              <span className="hidden md:inline">Live Video</span>
            </Button>
            <Button variant="ghost" size="sm" className="text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800">
              <Image className="w-5 h-5 mr-1 md:mr-2 text-green-500" />
              <span className="hidden md:inline">Photo/Video</span>
            </Button>
            <Button variant="ghost" size="sm" className="text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800">
              <Smile className="w-5 h-5 mr-1 md:mr-2 text-yellow-500" />
              <span className="hidden md:inline">Feeling/Activity</span>
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
