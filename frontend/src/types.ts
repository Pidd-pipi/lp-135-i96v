export type UserRole = 'user' | 'org' | 'admin';

export interface User {
  id: string;
  username: string;
  email: string;
  role: UserRole;
  avatar?: string;
  phone?: string;
  realName?: string;
  totalDonation: number;
  serviceHours: number;
  createdAt: string;
}

export type ProjectCategory = 'education' | 'elderly' | 'medical' | 'disaster' | 'environment' | 'other';
export type ProjectStatus = 'pending' | 'approved' | 'rejected' | 'completed';

export interface Project {
  id: string;
  organizationId: string;
  title: string;
  description?: string;
  category: ProjectCategory;
  targetAmount: number;
  currentAmount: number;
  executionPlan?: string;
  coverImage?: string;
  status: ProjectStatus;
  startDate?: string;
  endDate?: string;
  createdAt: string;
  progress?: number;
  organization?: Organization;
}

export type OrgStatus = 'pending' | 'approved' | 'rejected';

export interface Organization {
  id: string;
  userId: string;
  name: string;
  description?: string;
  licenseNumber?: string;
  contactPerson?: string;
  contactPhone?: string;
  address?: string;
  status: OrgStatus;
  createdAt: string;
}

export type PaymentMethod = 'wechat' | 'alipay' | 'bank';
export type PaymentStatus = 'pending' | 'success' | 'failed';

export interface Donation {
  id: string;
  userId: string;
  projectId: string;
  amount: number;
  paymentMethod: PaymentMethod;
  paymentStatus: PaymentStatus;
  transactionId?: string;
  certificateNo?: string;
  certificateUrl?: string;
  isAnonymous: boolean;
  message?: string;
  createdAt: string;
  donorName?: string;
  user?: User;
  project?: Project;
}

export interface ProjectUpdate {
  id: string;
  projectId: string;
  title: string;
  content?: string;
  images?: string;
  createdAt: string;
}

export interface RankingItem {
  id: string;
  userId: string;
  username: string;
  avatar?: string;
  realName?: string;
  totalDonation: number;
  serviceHours: number;
  rank: number;
}
