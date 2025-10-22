// User types
export interface User {
  id: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  bio?: string;
  avatarUrl?: string;
  location?: string;
  website?: string;
  role: 'user' | 'admin' | 'moderator';
  isVerified: boolean;
  isActive: boolean;
  preferences: UserPreferences;
  lastLoginAt?: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  notifications: {
    email: boolean;
    push: boolean;
    billUpdates: boolean;
    voteAlerts: boolean;
    politicianUpdates: boolean;
  };
  privacy: {
    profileVisibility: 'public' | 'private';
    showVotingHistory: boolean;
    showFollowing: boolean;
  };
}

// Politician types
export interface Politician {
  id: string;
  bioguideId?: string;
  govtrackId?: string;
  opensecretsId?: string;
  votesmartId?: string;
  fecIds?: string[];
  firstName: string;
  lastName: string;
  middleName?: string;
  suffix?: string;
  nickname?: string;
  fullName: string;
  gender?: 'M' | 'F' | 'Other';
  dateOfBirth?: Date;
  party: string;
  state: string;
  district?: string;
  chamber: 'house' | 'senate' | 'governor' | 'state_house' | 'state_senate';
  office: string;
  phone?: string;
  website?: string;
  contactForm?: string;
  twitterAccount?: string;
  facebookAccount?: string;
  youtubeAccount?: string;
  biography?: string;
  imageUrl?: string;
  inOffice: boolean;
  termStart?: Date;
  termEnd?: Date;
  termsServed: number;
  committees: Committee[];
  leadershipRoles: LeadershipRole[];
  influenceScore: number;
  bipartisanScore: number;
  attendanceRate: number;
  billsSponsored: number;
  billsCosponsored: number;
  policyPositions: Record<string, any>;
  votingPatterns: Record<string, any>;
  lastUpdated?: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface Committee {
  id: string;
  name: string;
  chamber: 'house' | 'senate' | 'joint';
  type: 'standing' | 'select' | 'special' | 'joint';
  role?: 'chair' | 'ranking_member' | 'member';
}

export interface LeadershipRole {
  title: string;
  startDate: Date;
  endDate?: Date;
  chamber: 'house' | 'senate';
}

// Bill types
export interface Bill {
  id: string;
  congressBillId: string;
  billNumber: string;
  congress: number;
  billType: 'hr' | 's' | 'hjres' | 'sjres' | 'hconres' | 'sconres' | 'hres' | 'sres';
  title: string;
  shortTitle?: string;
  officialTitle?: string;
  summary?: string;
  fullText?: string;
  sponsorId?: string;
  sponsor?: Politician;
  cosponsors: string[];
  status: BillStatus;
  actions: BillAction[];
  committees: string[];
  subjects: string[];
  policyAreas: string[];
  controversyScore: number;
  bipartisanSupport: number;
  totalVotes: number;
  yesVotes: number;
  noVotes: number;
  abstainVotes: number;
  introducedDate?: Date;
  lastActionDate?: Date;
  congressUrl?: string;
  govtrackUrl?: string;
  relatedBills: string[];
  nlpSummary?: string;
  keyTopics: string[];
  complexityScore: number;
  createdAt: Date;
  updatedAt: Date;
}

export type BillStatus = 
  | 'introduced'
  | 'referred'
  | 'reported'
  | 'passed_house'
  | 'passed_senate'
  | 'to_president'
  | 'signed'
  | 'vetoed'
  | 'failed';

export interface BillAction {
  date: Date;
  description: string;
  chamber?: 'house' | 'senate';
  actionType: string;
}

// Vote types
export interface Vote {
  id: string;
  voteId: string;
  billId?: string;
  bill?: Bill;
  politicianId: string;
  politician?: Politician;
  chamber: 'house' | 'senate';
  congress: number;
  session: number;
  rollCall: number;
  voteType: 'passage' | 'amendment' | 'procedural' | 'nomination' | 'other';
  votePosition: 'yes' | 'no' | 'present' | 'not_voting';
  question?: string;
  description?: string;
  result?: 'passed' | 'failed' | 'agreed_to' | 'rejected';
  totalYes: number;
  totalNo: number;
  totalPresent: number;
  totalNotVoting: number;
  voteDate: Date;
  voteTime?: string;
  congressUrl?: string;
  partyLineVote: boolean;
  bipartisanIndex: number;
  createdAt: Date;
  updatedAt: Date;
}

// Social types
export interface Post {
  id: string;
  userId: string;
  user?: User;
  politicianId?: string;
  politician?: Politician;
  billId?: string;
  bill?: Bill;
  content: string;
  mediaUrls: string[];
  postType: 'text' | 'image' | 'video' | 'poll' | 'bill_share' | 'vote_share';
  pollOptions?: PollOption[];
  likesCount: number;
  commentsCount: number;
  sharesCount: number;
  isPinned: boolean;
  isDeleted: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface PollOption {
  id: string;
  text: string;
  votes: number;
}

export interface Comment {
  id: string;
  postId: string;
  post?: Post;
  userId: string;
  user?: User;
  parentCommentId?: string;
  parentComment?: Comment;
  content: string;
  likesCount: number;
  repliesCount: number;
  isDeleted: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface Like {
  id: string;
  userId: string;
  user?: User;
  postId?: string;
  post?: Post;
  commentId?: string;
  comment?: Comment;
  likeType: 'like' | 'dislike' | 'love' | 'angry' | 'sad';
  createdAt: Date;
  updatedAt: Date;
}

export interface Follow {
  id: string;
  userId: string;
  user?: User;
  politicianId: string;
  politician?: Politician;
  notificationsEnabled: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserFollow {
  id: string;
  followerId: string;
  follower?: User;
  followingId: string;
  following?: User;
  createdAt: Date;
  updatedAt: Date;
}

export interface Notification {
  id: string;
  userId: string;
  user?: User;
  type: NotificationType;
  title: string;
  message?: string;
  data: Record<string, any>;
  isRead: boolean;
  createdAt: Date;
  updatedAt: Date;
}

export type NotificationType = 
  | 'like'
  | 'comment'
  | 'follow'
  | 'mention'
  | 'bill_update'
  | 'vote_alert'
  | 'politician_update';

// Analytics types
export interface PoliticianSimilarity {
  id: string;
  politicianAId: string;
  politicianA?: Politician;
  politicianBId: string;
  politicianB?: Politician;
  votingSimilarity: number;
  policySimilarity: number;
  overallSimilarity: number;
  commonVotes: number;
  agreementCount: number;
  similarityBreakdown: Record<string, any>;
  calculatedAt: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface BillTopic {
  id: string;
  billId: string;
  bill?: Bill;
  topic: string;
  confidence: number;
  extractionMethod?: string;
  keywords: string[];
  createdAt: Date;
  updatedAt: Date;
}

export interface VotingPattern {
  id: string;
  politicianId: string;
  politician?: Politician;
  topic: string;
  totalVotes: number;
  yesVotes: number;
  noVotes: number;
  abstainVotes: number;
  yesPercentage: number;
  partyAlignment: number;
  bipartisanScore: number;
  recentVotes: any[];
  lastCalculated: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface DailyAnalytics {
  id: string;
  date: Date;
  metricType: string;
  data: Record<string, any>;
  createdAt: Date;
  updatedAt: Date;
}

export interface UserEngagement {
  id: string;
  userId: string;
  user?: User;
  date: Date;
  postsCreated: number;
  commentsMade: number;
  likesGiven: number;
  politiciansFollowed: number;
  billsViewed: number;
  votesAnalyzed: number;
  sessionDuration: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface TrendingTopic {
  id: string;
  topic: string;
  mentionCount: number;
  engagementScore: number;
  growthRate: number;
  relatedBills: string[];
  relatedPoliticians: string[];
  trendingDate: Date;
  createdAt: Date;
  updatedAt: Date;
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  pagination?: PaginationInfo;
}

export interface PaginationInfo {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

// Authentication types
export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  firstName: string;
  lastName: string;
}

// Search and filter types
export interface SearchFilters {
  query?: string;
  party?: string;
  state?: string;
  chamber?: string;
  status?: string;
  dateFrom?: Date;
  dateTo?: Date;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

// External API types
export interface CongressApiResponse {
  bills?: any[];
  members?: any[];
  votes?: any[];
  pagination?: {
    count: number;
    next?: string;
    previous?: string;
  };
}

export interface OpenStatesApiResponse {
  results?: any[];
  meta?: {
    page: number;
    max_page: number;
    per_page: number;
    total_count: number;
  };
}

