import React, { useState } from 'react';
import Head from 'next/head';
import { Search, UserPlus, UserCheck } from 'lucide-react';
import { Navbar } from '@/components/facebook/Navbar';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';

const FriendsPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'suggestions' | 'requests' | 'all'>('suggestions');

  const friendSuggestions = [
    { id: '1', name: 'Alex Thompson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Alex', mutualFriends: 12 },
    { id: '2', name: 'Jessica Lee', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Jessica', mutualFriends: 8 },
    { id: '3', name: 'Ryan Parker', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Ryan', mutualFriends: 15 },
    { id: '4', name: 'Sophia Martinez', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sophia', mutualFriends: 6 },
    { id: '5', name: 'Daniel White', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Daniel', mutualFriends: 20 },
    { id: '6', name: 'Mia Johnson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Mia', mutualFriends: 4 },
  ];

  const friendRequests = [
    { id: '1', name: 'Chris Anderson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Chris', mutualFriends: 5, time: '2 hours ago' },
    { id: '2', name: 'Lauren Smith', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Lauren', mutualFriends: 3, time: '1 day ago' },
  ];

  const allFriends = [
    { id: '1', name: 'Sarah Wilson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah' },
    { id: '2', name: 'Michael Chen', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Michael' },
    { id: '3', name: 'Emma Davis', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Emma' },
    { id: '4', name: 'James Rodriguez', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=James' },
    { id: '5', name: 'Olivia Taylor', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Olivia' },
    { id: '6', name: 'William Brown', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=William' },
    { id: '7', name: 'Ava Garcia', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Ava' },
    { id: '8', name: 'Liam Wilson', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Liam' },
  ];

  return (
    <>
      <Head>
        <title>Friends - Facebook Clone</title>
        <meta name="description" content="Find and connect with friends" />
      </Head>

      <div className="min-h-screen bg-gray-100 dark:bg-gray-950">
        <Navbar />

        <div className="max-w-7xl mx-auto p-4">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
            {/* Left Sidebar */}
            <Card className="lg:col-span-1 h-fit">
              <CardHeader>
                <CardTitle>Friends</CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                <div className="space-y-1">
                  <Button
                    variant="ghost"
                    className={`w-full justify-start px-4 py-3 ${
                      activeTab === 'suggestions' ? 'bg-blue-50 dark:bg-blue-950 text-blue-600' : ''
                    }`}
                    onClick={() => setActiveTab('suggestions')}
                  >
                    <UserPlus className="w-5 h-5 mr-3" />
                    Friend Suggestions
                  </Button>
                  <Button
                    variant="ghost"
                    className={`w-full justify-start px-4 py-3 ${
                      activeTab === 'requests' ? 'bg-blue-50 dark:bg-blue-950 text-blue-600' : ''
                    }`}
                    onClick={() => setActiveTab('requests')}
                  >
                    <UserCheck className="w-5 h-5 mr-3" />
                    Friend Requests
                    {friendRequests.length > 0 && (
                      <span className="ml-auto bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                        {friendRequests.length}
                      </span>
                    )}
                  </Button>
                  <Button
                    variant="ghost"
                    className={`w-full justify-start px-4 py-3 ${
                      activeTab === 'all' ? 'bg-blue-50 dark:bg-blue-950 text-blue-600' : ''
                    }`}
                    onClick={() => setActiveTab('all')}
                  >
                    All Friends
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Main Content */}
            <div className="lg:col-span-3">
              {/* Search Bar */}
              <Card className="mb-4">
                <CardContent className="p-4">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <Input
                      type="text"
                      placeholder="Search friends"
                      className="pl-10 bg-gray-100 dark:bg-gray-800 border-none"
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Friend Suggestions */}
              {activeTab === 'suggestions' && (
                <>
                  <h2 className="text-xl font-bold mb-4">People You May Know</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                    {friendSuggestions.map((person) => (
                      <Card key={person.id}>
                        <CardContent className="p-0">
                          <div className="aspect-square bg-gradient-to-br from-blue-100 to-purple-100 dark:from-blue-900 dark:to-purple-900 flex items-center justify-center">
                            <Avatar className="w-32 h-32">
                              <AvatarImage src={person.avatar} />
                              <AvatarFallback className="text-4xl">{person.name[0]}</AvatarFallback>
                            </Avatar>
                          </div>
                          <div className="p-4">
                            <h3 className="font-semibold text-lg mb-1">{person.name}</h3>
                            <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                              {person.mutualFriends} mutual friends
                            </p>
                            <div className="space-y-2">
                              <Button variant="default" className="w-full">
                                Add Friend
                              </Button>
                              <Button variant="secondary" className="w-full">
                                Remove
                              </Button>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </>
              )}

              {/* Friend Requests */}
              {activeTab === 'requests' && (
                <>
                  <h2 className="text-xl font-bold mb-4">Friend Requests</h2>
                  {friendRequests.length > 0 ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                      {friendRequests.map((person) => (
                        <Card key={person.id}>
                          <CardContent className="p-0">
                            <div className="aspect-square bg-gradient-to-br from-green-100 to-blue-100 dark:from-green-900 dark:to-blue-900 flex items-center justify-center">
                              <Avatar className="w-32 h-32">
                                <AvatarImage src={person.avatar} />
                                <AvatarFallback className="text-4xl">{person.name[0]}</AvatarFallback>
                              </Avatar>
                            </div>
                            <div className="p-4">
                              <h3 className="font-semibold text-lg mb-1">{person.name}</h3>
                              <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                                {person.mutualFriends} mutual friends
                              </p>
                              <p className="text-xs text-gray-500 dark:text-gray-500 mb-3">
                                {person.time}
                              </p>
                              <div className="space-y-2">
                                <Button variant="default" className="w-full">
                                  Confirm
                                </Button>
                                <Button variant="secondary" className="w-full">
                                  Delete
                                </Button>
                              </div>
                            </div>
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  ) : (
                    <Card>
                      <CardContent className="p-8 text-center">
                        <p className="text-gray-600 dark:text-gray-400">
                          No friend requests at the moment
                        </p>
                      </CardContent>
                    </Card>
                  )}
                </>
              )}

              {/* All Friends */}
              {activeTab === 'all' && (
                <>
                  <h2 className="text-xl font-bold mb-4">All Friends ({allFriends.length})</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                    {allFriends.map((friend) => (
                      <Card key={friend.id}>
                        <CardContent className="p-4 flex items-center space-x-3">
                          <Avatar className="w-16 h-16">
                            <AvatarImage src={friend.avatar} />
                            <AvatarFallback className="text-xl">{friend.name[0]}</AvatarFallback>
                          </Avatar>
                          <div className="flex-1">
                            <h3 className="font-semibold">{friend.name}</h3>
                            <Button variant="link" className="p-0 h-auto text-sm text-gray-600 dark:text-gray-400">
                              Message
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default FriendsPage;
