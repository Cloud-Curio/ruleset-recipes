import React from 'react';
import Head from 'next/head';
import { Navbar } from '@/components/facebook/Navbar';
import { LeftSidebar } from '@/components/facebook/LeftSidebar';
import { RightSidebar } from '@/components/facebook/RightSidebar';
import { Stories } from '@/components/facebook/Stories';
import { PostComposer } from '@/components/facebook/PostComposer';
import { Post } from '@/components/facebook/Post';

const FeedPage: React.FC = () => {
  // Mock posts data
  const posts = [
    {
      id: '1',
      author: {
        name: 'Sarah Wilson',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah',
      },
      timestamp: '2 hours ago',
      content: 'Just finished an amazing project using React and Next.js! The developer experience is incredible. 🚀',
      image: 'https://picsum.photos/seed/react/800/600',
      likes: 124,
      comments: 23,
      shares: 5,
    },
    {
      id: '2',
      author: {
        name: 'Michael Chen',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Michael',
      },
      timestamp: '5 hours ago',
      content: 'Beautiful sunset at the beach today! Nature never fails to amaze me. 🌅',
      image: 'https://picsum.photos/seed/sunset/800/600',
      likes: 342,
      comments: 67,
      shares: 12,
    },
    {
      id: '3',
      author: {
        name: 'Emma Davis',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma',
      },
      timestamp: '8 hours ago',
      content: 'Had the most productive day! Completed all my tasks and still had time for a workout. Feeling accomplished! 💪',
      likes: 89,
      comments: 15,
      shares: 3,
    },
    {
      id: '4',
      author: {
        name: 'James Rodriguez',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=James',
      },
      timestamp: '12 hours ago',
      content: 'Excited to share that I just got promoted! Hard work really does pay off. Thank you to everyone who supported me on this journey! 🎉',
      image: 'https://picsum.photos/seed/celebration/800/600',
      likes: 567,
      comments: 145,
      shares: 23,
    },
    {
      id: '5',
      author: {
        name: 'Olivia Taylor',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Olivia',
      },
      timestamp: '1 day ago',
      content: 'Coffee and code - the perfect combination for a Saturday morning! ☕️💻',
      image: 'https://picsum.photos/seed/coffee/800/600',
      likes: 201,
      comments: 34,
      shares: 8,
    },
  ];

  return (
    <>
      <Head>
        <title>Facebook Clone - News Feed</title>
        <meta name="description" content="Stay connected with friends and the world around you on Facebook." />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>

      <div className="min-h-screen bg-gray-100 dark:bg-gray-950">
        <Navbar />
        
        <div className="flex max-w-full">
          {/* Left Sidebar */}
          <LeftSidebar />

          {/* Main Feed */}
          <main className="flex-1 max-w-2xl mx-auto px-4 py-4">
            <Stories />
            <PostComposer />
            
            {/* Posts */}
            <div className="space-y-4">
              {posts.map((post) => (
                <Post key={post.id} {...post} />
              ))}
            </div>

            {/* Loading Indicator */}
            <div className="text-center py-8">
              <p className="text-gray-500 dark:text-gray-400">
                You're all caught up!
              </p>
            </div>
          </main>

          {/* Right Sidebar */}
          <RightSidebar />
        </div>
      </div>
    </>
  );
};

export default FeedPage;
