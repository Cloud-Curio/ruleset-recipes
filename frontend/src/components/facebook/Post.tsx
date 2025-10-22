import React, { useState } from 'react';
import { 
  MoreHorizontal, 
  ThumbsUp, 
  MessageCircle, 
  Share2, 
  X,
  Heart,
  Laugh,
  Angry
} from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/contexts/AuthContext';

interface PostProps {
  id: string;
  author: {
    name: string;
    avatar: string;
  };
  timestamp: string;
  content: string;
  image?: string;
  likes: number;
  comments: number;
  shares: number;
}

export const Post: React.FC<PostProps> = ({
  author,
  timestamp,
  content,
  image,
  likes,
  comments,
  shares,
}) => {
  const { user } = useAuth();
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(likes);
  const [showComments, setShowComments] = useState(false);
  const [commentText, setCommentText] = useState('');

  const handleLike = () => {
    setLiked(!liked);
    setLikeCount(liked ? likeCount - 1 : likeCount + 1);
  };

  const handleComment = () => {
    if (commentText.trim()) {
      console.log('Comment:', commentText);
      setCommentText('');
    }
  };

  return (
    <Card className="mb-4">
      <CardContent className="p-0">
        {/* Post Header */}
        <div className="p-4 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Avatar className="w-10 h-10">
              <AvatarImage src={author.avatar} />
              <AvatarFallback>{author.name[0]}</AvatarFallback>
            </Avatar>
            <div>
              <p className="font-semibold text-sm">{author.name}</p>
              <p className="text-xs text-gray-500 dark:text-gray-400">{timestamp}</p>
            </div>
          </div>
          <Button variant="ghost" size="icon" className="rounded-full">
            <MoreHorizontal className="w-5 h-5" />
          </Button>
        </div>

        {/* Post Content */}
        <div className="px-4 pb-3">
          <p className="text-sm whitespace-pre-wrap">{content}</p>
        </div>

        {/* Post Image */}
        {image && (
          <div className="relative">
            <img 
              src={image} 
              alt="Post content" 
              className="w-full object-cover max-h-[600px]"
            />
          </div>
        )}

        {/* Reactions Summary */}
        <div className="px-4 py-2 flex items-center justify-between text-sm text-gray-600 dark:text-gray-400">
          <div className="flex items-center space-x-1">
            <div className="flex -space-x-1">
              <div className="w-5 h-5 bg-blue-600 rounded-full flex items-center justify-center">
                <ThumbsUp className="w-3 h-3 text-white fill-white" />
              </div>
              <div className="w-5 h-5 bg-red-600 rounded-full flex items-center justify-center">
                <Heart className="w-3 h-3 text-white fill-white" />
              </div>
            </div>
            <span className="hover:underline cursor-pointer">{likeCount}</span>
          </div>
          <div className="flex space-x-3">
            <span className="hover:underline cursor-pointer">{comments} comments</span>
            <span className="hover:underline cursor-pointer">{shares} shares</span>
          </div>
        </div>

        <Separator />

        {/* Action Buttons */}
        <div className="px-2 py-1 flex justify-around">
          <Button 
            variant="ghost" 
            className={`flex-1 rounded-lg ${liked ? 'text-blue-600' : 'text-gray-600 dark:text-gray-400'}`}
            onClick={handleLike}
          >
            <ThumbsUp className={`w-5 h-5 mr-2 ${liked ? 'fill-blue-600' : ''}`} />
            <span className="font-semibold">Like</span>
          </Button>
          <Button 
            variant="ghost" 
            className="flex-1 text-gray-600 dark:text-gray-400 rounded-lg"
            onClick={() => setShowComments(!showComments)}
          >
            <MessageCircle className="w-5 h-5 mr-2" />
            <span className="font-semibold">Comment</span>
          </Button>
          <Button variant="ghost" className="flex-1 text-gray-600 dark:text-gray-400 rounded-lg">
            <Share2 className="w-5 h-5 mr-2" />
            <span className="font-semibold">Share</span>
          </Button>
        </div>

        {/* Comments Section */}
        {showComments && (
          <>
            <Separator />
            <div className="p-4">
              <div className="flex space-x-2">
                <Avatar className="w-8 h-8">
                  <AvatarImage src={user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=User'} />
                  <AvatarFallback>{user?.name?.[0] || 'U'}</AvatarFallback>
                </Avatar>
                <div className="flex-1 flex space-x-2">
                  <Input
                    placeholder="Write a comment..."
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && handleComment()}
                    className="flex-1 rounded-full bg-gray-100 dark:bg-gray-800 border-none"
                  />
                </div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
};
