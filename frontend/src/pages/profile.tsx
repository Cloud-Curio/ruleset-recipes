import React from 'react';
import Head from 'next/head';
import { Camera, MapPin, Briefcase, GraduationCap, Heart, MoreHorizontal, Edit } from 'lucide-react';
import { Navbar } from '@/components/facebook/Navbar';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Post } from '@/components/facebook/Post';
import { useAuth } from '@/contexts/AuthContext';

const ProfilePage: React.FC = () => {
  const { user } = useAuth();

  const profileData = {
    name: user?.name || 'John Doe',
    avatar: user?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=John',
    coverImage: 'https://picsum.photos/seed/cover/1200/400',
    bio: 'Software Developer | React Enthusiast | Coffee Lover ☕',
    location: 'San Francisco, CA',
    work: 'Software Engineer at Tech Corp',
    education: 'Stanford University',
    relationshipStatus: 'Single',
    friends: 482,
    photos: 127,
  };

  const posts = [
    {
      id: '1',
      author: {
        name: profileData.name,
        avatar: profileData.avatar,
      },
      timestamp: '3 hours ago',
      content: 'Just deployed a new feature! Feeling accomplished 🚀',
      image: 'https://picsum.photos/seed/feature/800/600',
      likes: 45,
      comments: 12,
      shares: 2,
    },
    {
      id: '2',
      author: {
        name: profileData.name,
        avatar: profileData.avatar,
      },
      timestamp: '1 day ago',
      content: 'Weekend vibes! Taking a break from coding to enjoy nature.',
      image: 'https://picsum.photos/seed/nature/800/600',
      likes: 78,
      comments: 23,
      shares: 5,
    },
  ];

  const photos = Array.from({ length: 9 }, (_, i) => ({
    id: i + 1,
    url: `https://picsum.photos/seed/photo${i}/300/300`,
  }));

  const friends = [
    { id: '1', name: 'Sarah Wilson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah' },
    { id: '2', name: 'Michael Chen', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Michael' },
    { id: '3', name: 'Emma Davis', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma' },
    { id: '4', name: 'James Rodriguez', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=James' },
    { id: '5', name: 'Olivia Taylor', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Olivia' },
    { id: '6', name: 'William Brown', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=William' },
  ];

  return (
    <>
      <Head>
        <title>{profileData.name} - Profile</title>
        <meta name="description" content={`${profileData.name}'s profile on Facebook`} />
      </Head>

      <div className="min-h-screen bg-gray-100 dark:bg-gray-950">
        <Navbar />

        <div className="max-w-5xl mx-auto">
          {/* Cover Photo */}
          <div className="relative">
            <div className="h-96 bg-gradient-to-r from-blue-400 to-purple-500 relative overflow-hidden">
              <img 
                src={profileData.coverImage} 
                alt="Cover" 
                className="w-full h-full object-cover"
              />
              <Button 
                variant="secondary" 
                className="absolute bottom-4 right-4"
              >
                <Camera className="w-4 h-4 mr-2" />
                Edit Cover Photo
              </Button>
            </div>

            {/* Profile Info Header */}
            <div className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800">
              <div className="px-4 pb-4">
                <div className="flex flex-col md:flex-row md:items-end md:justify-between">
                  <div className="flex flex-col md:flex-row md:items-end -mt-20 md:-mt-24">
                    <div className="relative">
                      <Avatar className="w-40 h-40 border-4 border-white dark:border-gray-900">
                        <AvatarImage src={profileData.avatar} />
                        <AvatarFallback className="text-4xl">{profileData.name[0]}</AvatarFallback>
                      </Avatar>
                      <Button 
                        size="icon" 
                        variant="secondary" 
                        className="absolute bottom-2 right-2 rounded-full w-9 h-9"
                      >
                        <Camera className="w-4 h-4" />
                      </Button>
                    </div>
                    <div className="mt-4 md:mt-0 md:ml-6 md:mb-4">
                      <h1 className="text-3xl font-bold">{profileData.name}</h1>
                      <p className="text-gray-600 dark:text-gray-400">{profileData.friends} friends</p>
                    </div>
                  </div>
                  <div className="flex space-x-2 mt-4 md:mt-0 md:mb-4">
                    <Button variant="default">
                      <Edit className="w-4 h-4 mr-2" />
                      Edit Profile
                    </Button>
                    <Button variant="secondary">
                      <MoreHorizontal className="w-4 h-4" />
                    </Button>
                  </div>
                </div>

                {/* Navigation Tabs */}
                <div className="flex space-x-2 mt-4 overflow-x-auto">
                  <Button variant="ghost" className="rounded-lg border-b-2 border-blue-600">
                    Posts
                  </Button>
                  <Button variant="ghost" className="rounded-lg">
                    About
                  </Button>
                  <Button variant="ghost" className="rounded-lg">
                    Friends
                  </Button>
                  <Button variant="ghost" className="rounded-lg">
                    Photos
                  </Button>
                  <Button variant="ghost" className="rounded-lg">
                    Videos
                  </Button>
                  <Button variant="ghost" className="rounded-lg">
                    More
                  </Button>
                </div>
              </div>
            </div>
          </div>

          {/* Profile Content */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 p-4">
            {/* Left Column - About */}
            <div className="lg:col-span-1 space-y-4">
              {/* Intro Card */}
              <Card>
                <CardContent className="p-4">
                  <h2 className="text-xl font-bold mb-4">Intro</h2>
                  <p className="text-center text-gray-600 dark:text-gray-400 mb-4">
                    {profileData.bio}
                  </p>
                  <Button variant="secondary" className="w-full mb-4">
                    Edit Bio
                  </Button>
                  <div className="space-y-3">
                    <div className="flex items-center text-sm">
                      <Briefcase className="w-4 h-4 mr-3 text-gray-500" />
                      <span>{profileData.work}</span>
                    </div>
                    <div className="flex items-center text-sm">
                      <GraduationCap className="w-4 h-4 mr-3 text-gray-500" />
                      <span>Studied at {profileData.education}</span>
                    </div>
                    <div className="flex items-center text-sm">
                      <MapPin className="w-4 h-4 mr-3 text-gray-500" />
                      <span>Lives in {profileData.location}</span>
                    </div>
                    <div className="flex items-center text-sm">
                      <Heart className="w-4 h-4 mr-3 text-gray-500" />
                      <span>{profileData.relationshipStatus}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Photos Card */}
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between mb-4">
                    <h2 className="text-xl font-bold">Photos</h2>
                    <Button variant="link" className="text-blue-600">
                      See all photos
                    </Button>
                  </div>
                  <div className="grid grid-cols-3 gap-2">
                    {photos.map((photo) => (
                      <div key={photo.id} className="aspect-square overflow-hidden rounded-lg">
                        <img 
                          src={photo.url} 
                          alt={`Photo ${photo.id}`}
                          className="w-full h-full object-cover hover:scale-110 transition-transform duration-200 cursor-pointer"
                        />
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              {/* Friends Card */}
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h2 className="text-xl font-bold">Friends</h2>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {profileData.friends} friends
                      </p>
                    </div>
                    <Button variant="link" className="text-blue-600">
                      See all
                    </Button>
                  </div>
                  <div className="grid grid-cols-3 gap-2">
                    {friends.map((friend) => (
                      <div key={friend.id} className="flex flex-col items-center">
                        <Avatar className="w-20 h-20 mb-2">
                          <AvatarImage src={friend.avatar} />
                          <AvatarFallback>{friend.name[0]}</AvatarFallback>
                        </Avatar>
                        <p className="text-xs font-medium text-center line-clamp-2">
                          {friend.name}
                        </p>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Right Column - Posts */}
            <div className="lg:col-span-2 space-y-4">
              {posts.map((post) => (
                <Post key={post.id} {...post} />
              ))}
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default ProfilePage;
